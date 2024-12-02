package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *slog.Logger
)

func Init() {
	r := &lumberjack.Logger{
		Filename:   fmt.Sprintf(global.DataDir+"logs/%s.log", "monitor"),
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	log := slog.NewJSONHandler(r, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(log)
}

// Control monitor task
func Control(ctx context.Context) (any, error) {
	data, err := GetMonitorData()
	if err != nil {
		return nil, err
	}
	logger.Debug("Monitor",
		slog.Float64("CpuPercent", data.CpuPercent),
		slog.Float64("MemPercent", data.MemPercent),
		slog.Int("ThreadCount", data.ThreadCount),
		slog.Int("GoroutineCount", data.GoroutineCount),
		slog.Float64("IOPercent", data.IOPercent),
		slog.Float64("Load1", data.Load1),
		slog.Float64("Load5", data.Load5),
		slog.Float64("Load15", data.Load15),
		slog.Any("NetStats", data.NetStats),
		slog.Int64("Timestamp", data.Timestamp),
	)
	return true, nil
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

type SystemInfo struct {
	OsName        string  `json:"osName"`        // 操作系统名称
	KernelVersion string  `json:"kernelVersion"` // 内核版本
	MemoryTotal   string  `json:"memoryTotal"`   // 内存总量
	DiskTotal     string  `json:"diskTotal"`     // 硬盘总量
	CpuArch       string  `json:"cpuArch"`       // CPU架构
	CpuCount      int     `json:"cpuCount"`      // CPU数量
	CpuModel      string  `json:"cpuModel"`      // CPU型号
	CpuPercent    float64 `json:"cpuPercent"`    // CPU使用率
}

func System(ctx context.Context) *SystemInfo {
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// 计算所有磁盘总容量
	partitions, _ := disk.Partitions(false)
	var totalSize uint64
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		totalSize += usage.Total
	}

	return &SystemInfo{
		OsName:        hostInfo.Platform + " " + hostInfo.PlatformVersion,
		KernelVersion: hostInfo.KernelVersion,
		MemoryTotal:   humanize.Bytes(memInfo.Total),
		DiskTotal:     humanize.Bytes(totalSize),
		CpuArch:       runtime.GOARCH,
		CpuCount:      runtime.NumCPU(),
		CpuModel:      cpuInfo[0].ModelName,
		CpuPercent:    helper.Round(cpuPercent[0], 2),
	}
}
