package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/samber/do/v2"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/spf13/afero"
)

type Monitor struct {
	OsName      string // OS
	BootTime    string // Startup time
	LogSize     uint64 // The size of the log directory
	LogSizeF    string
	UploadSize  uint64 // The size of the upload directory
	UploadSizeF string
	TmpSize     uint64 // The size of the cache directory.
	TmpSizeF    string
}

// GetMonitorInfo Retrieve monitoring information
func GetMonitorInfo() *Monitor {
	data := Monitor{}
	data.LogSize = getDirSize("/data/logs")
	data.LogSizeF = humanize.Bytes(data.LogSize)
	data.UploadSize = getDirSize("/public/uploads")
	data.UploadSizeF = humanize.Bytes(data.UploadSize)
	data.TmpSize = getDirSize("/tmp")
	data.TmpSizeF = humanize.Bytes(data.TmpSize)
	data.BootTime = global.BootTime.Format("2006-01-02 15:04:05")
	sysInfo, _ := host.Info()
	data.OsName = sysInfo.Platform + " " + sysInfo.PlatformVersion
	return &data

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

	// 获取内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %w", err)
	}

	threadCount := pprof.Lookup("threadcreate").Count()
	GoroutineCount := runtime.NumGoroutine()

	// 获取磁盘IO统计
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

	// 获取系统负载
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get system load: %w", err)
	}

	// 获取网络速率
	prevNet, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network statistics: %w", err)
	}
	time.Sleep(time.Second)
	currentNet, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network statistics: %w", err)
	}

	// 计算每个网卡的网络速率
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
		CpuPercent:     helper.Round(cpuPercent[0], 2), // cpuPercent返回的是切片
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
