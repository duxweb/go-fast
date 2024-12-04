package monitor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/spf13/afero"
)

type DirSize struct {
	Path  string `json:"path"`
	Size  uint64 `json:"size"`
	SizeF string `json:"sizeF"`
}

type Monitor struct {
	OsName   string    `json:"osName"`   // OS 名称
	BootTime string    `json:"bootTime"` // 启动时间
	DirSize  []DirSize `json:"dirSize"`  // 目录大小

	KernelVersion string `json:"kernelVersion"` // 内核版本
	MemoryTotal   string `json:"memoryTotal"`   // 内存总量
	DiskTotal     string `json:"diskTotal"`     // 硬盘总量
	CpuArch       string `json:"cpuArch"`       // CPU架构
	CpuCount      int    `json:"cpuCount"`      // CPU数量
	CpuModel      string `json:"cpuModel"`      // CPU型号

	DefaultDatabase string `json:"defaultDatabase"` // 默认数据库
	DefaultQueue    string `json:"defaultQueue"`    // 默认队列
	Debug           bool   `json:"debug"`           // 是否是debug模式
	DefaultLang     string `json:"defaultLang"`     // 默认语言
	DefaultTimezone string `json:"defaultTimezone"` // 默认时区

	Port            string `json:"port"`            // 端口
	Version         string `json:"version"`         // 版本
	GoVersion       string `json:"goVersion"`       // Go版本
	DatabaseVersion string `json:"databaseVersion"` // 数据库版本
}

// GetMonitorInfo Retrieve monitoring information
func GetMonitorInfo() *Monitor {
	data := Monitor{}

	logSize := getDirSize("/data/logs")
	logSizeF := humanize.Bytes(logSize)
	uploadSize := getDirSize("/public/uploads")
	uploadSizeF := humanize.Bytes(uploadSize)
	tmpSize := getDirSize("/tmp")
	tmpSizeF := humanize.Bytes(tmpSize)

	data.DirSize = []DirSize{
		{Path: "/data/logs", Size: logSize, SizeF: logSizeF},
		{Path: "/public/uploads", Size: uploadSize, SizeF: uploadSizeF},
		{Path: "/tmp", Size: tmpSize, SizeF: tmpSizeF},
	}

	data.BootTime = global.BootTime.Format("2006-01-02 15:04:05")
	sysInfo, _ := host.Info()
	data.OsName = sysInfo.Platform + " " + sysInfo.PlatformVersion

	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()

	diskStats, _ := disk.IOCounters()
	var totalSize uint64
	for name := range diskStats {
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
			continue
		}

		// 获取分区信息
		partitions, err := disk.Partitions(true)
		if err != nil {
			continue
		}

		// 查找对应分区的容量
		for _, partition := range partitions {
			if strings.Contains(partition.Device, name) {
				usage, err := disk.Usage(partition.Mountpoint)
				if err != nil {
					continue
				}
				totalSize += usage.Total
				break
			}
		}
	}

	data.KernelVersion = sysInfo.KernelVersion
	data.MemoryTotal = humanize.IBytes(memInfo.Total)
	data.DiskTotal = humanize.IBytes(totalSize)
	data.CpuArch = runtime.GOARCH
	data.CpuCount = runtime.NumCPU()
	data.CpuModel = cpuInfo[0].ModelName

	data.DefaultDatabase = config.Load("database").GetString("db.drivers.default.type")
	data.DefaultQueue = lo.Ternary(config.Load("use").GetString("queue.driver") == "redis", "Redis", "Sqlite")
	data.Debug = global.Debug
	data.DefaultLang = global.Lang
	data.DefaultTimezone = global.TimeLocation.String()
	data.Port = config.Load("use").GetString("server.port")
	data.Version = global.Version
	data.GoVersion = runtime.Version()

	//
	switch data.DefaultDatabase {
	case "mysql":
		var version string
		if err := database.Gorm().Raw("SELECT VERSION()").Scan(&version).Error; err == nil {
			data.DatabaseVersion = version
		}
		data.DefaultDatabase = "MySQL"
	case "sqlite":
		var version string
		if err := database.Gorm().Raw("SELECT sqlite_version()").Scan(&version).Error; err == nil {
			data.DatabaseVersion = version
		}
		data.DefaultDatabase = "SQLite"
	case "postgres":
		var version string
		if err := database.Gorm().Raw("SHOW server_version").Scan(&version).Error; err == nil {
			data.DatabaseVersion = version
		}
		data.DefaultDatabase = "PostgreSQL"
	}

	return &data
}

type DiskStatus struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

func Disk(ctx context.Context) []DiskStatus {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil
	}

	var diskStats []DiskStatus
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskStats = append(diskStats, DiskStatus{
			Path:        partition.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: helper.Round(usage.UsedPercent, 2),
		})
	}

	return diskStats
}

type NetStats struct {
	Name         string  `json:"name"`         // 网卡名称
	UploadRate   float64 `json:"uploadRate"`   // 上传速率 Mbps
	DownloadRate float64 `json:"downloadRate"` // 下载速率 Mbps
}

type MonitorData struct {
	CpuPercent     float64
	MemPercent     float64
	ThreadCount    int
	GoroutineCount int
	Timestamp      int64
	// IO 负载相关
	IOPercent float64 // IO 使用率（百分比）
	// CPU 负载
	Load1  float64 // 1分钟负载
	Load5  float64 // 5分钟负载
	Load15 float64 // 15分钟负载
	// 网络相关
	NetStats []NetStats
}

// GetMonitorData Retrieve monitoring data
func GetMonitorData() (*MonitorData, error) {
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %w", err)
	}

	threadCount := pprof.Lookup("threadcreate").Count()
	GoroutineCount := runtime.NumGoroutine()

	prevIO, err := disk.IOCounters()
	ioUsage := 0.0
	if err == nil {
		time.Sleep(time.Second) // 1秒采样间隔
		currentIO, err := disk.IOCounters()
		if err == nil {
			var totalIOTime float64
			for name, current := range currentIO {
				if prev, ok := prevIO[name]; ok {
					// 计算1秒内的IO时间变化
					deltaIOTime := float64(current.IoTime - prev.IoTime)
					// 转换为百分比 (考虑多核CPU)
					totalIOTime += deltaIOTime / (1000 * float64(runtime.NumCPU())) * 100
				}
			}
			ioUsage = helper.Round(totalIOTime, 2)
		}
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get system load: %w", err)
	}

	prevNet, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network statistics: %w", err)
	}
	time.Sleep(time.Second)
	currentNet, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network statistics: %w", err)
	}

	netStats := make([]NetStats, 0)
	for i, current := range currentNet {
		if i < len(prevNet) {
			prev := prevNet[i]
			uploadRate := float64(current.BytesSent-prev.BytesSent) * 8 / 1024 / 1024
			downloadRate := float64(current.BytesRecv-prev.BytesRecv) * 8 / 1024 / 1024
			netStats = append(netStats, NetStats{
				Name:         current.Name,
				UploadRate:   helper.Round(uploadRate, 2),
				DownloadRate: helper.Round(downloadRate, 2),
			})
		}
	}

	return &MonitorData{
		CpuPercent:     helper.Round(cpuPercent[0], 2),
		MemPercent:     helper.Round(memInfo.UsedPercent, 2),
		ThreadCount:    threadCount,
		GoroutineCount: GoroutineCount,
		Timestamp:      time.Now().UnixMilli(),

		IOPercent: helper.Round(ioUsage, 2),

		Load1:  helper.Round(loadAvg.Load1, 2),
		Load5:  helper.Round(loadAvg.Load5, 2),
		Load15: helper.Round(loadAvg.Load15, 2),

		NetStats: netStats,
	}, nil
}

// GetMonitorLog Retrieve monitoring logs
func GetMonitorLog() []map[string]any {
	loadFiles, _ := filepath.Glob(global.DataDir + "logs/monitor*.log")
	loadData := passingFiles(loadFiles)
	return loadData
}

func getDirSize(path string) uint64 {
	var size int64
	wd, _ := os.Getwd()

	ofs := do.MustInvokeNamed[afero.Fs](global.Injector, "os.fs")

	exists, _ := afero.Exists(ofs, wd+path)
	if !exists {
		return 0
	}
	_ = afero.Walk(ofs, wd+path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return uint64(size)
}

func passingFiles(files []string) []map[string]any {
	loadData := make([]map[string]any, 0)
	for _, file := range files {
		fileData, err := parsingFile(file)
		if err != nil {
			continue
		}
		loadData = append(loadData, fileData...)
	}
	return loadData
}

func parsingFile(file string) ([]map[string]any, error) {

	ofs := do.MustInvokeNamed[afero.Fs](global.Injector, "os.fs")
	fd, err := ofs.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	// 设置更大的buffer以提高性能
	const maxCapacity = 512 * 1024 // 512KB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	data := make([]map[string]any, 0, 100) // 预分配容量
	for scanner.Scan() {
		curData := map[string]any{}
		if err := json.Unmarshal(scanner.Bytes(), &curData); err != nil {
			continue
		}
		data = append(data, curData)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return data, nil
}
