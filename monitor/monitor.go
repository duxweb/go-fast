package monitor

import (
	"bufio"
	"encoding/json"
	"github.com/dustin/go-humanize"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

type MonitorInfo struct {
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
func GetMonitorInfo() *MonitorInfo {
	data := MonitorInfo{}
	data.LogSize = getDirSize("/logs")
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

type MonitorData struct {
	CpuPercent     float64
	MemPercent     float64
	ThreadCount    int
	GoroutineCount int
	Timestamp      int64
}

// GetMonitorData Retrieve monitoring data
func GetMonitorData() *MonitorData {
	p, _ := process.NewProcess(int32(os.Getpid()))
	cpuPercent, _ := p.Percent(time.Second)
	memPercent, _ := p.MemoryPercent()
	threadCount := pprof.Lookup("threadcreate").Count()
	GoroutineCount := runtime.NumGoroutine()

	return &MonitorData{
		CpuPercent:     helper.Round(cpuPercent, 2),
		MemPercent:     helper.Round(float64(memPercent), 2),
		ThreadCount:    threadCount,
		GoroutineCount: GoroutineCount,
		Timestamp:      time.Now().UnixMilli(),
	}
}

// GetMonitorLog Retrieve monitoring logs
func GetMonitorLog() []map[string]any {
	path := config.Get("app").GetString("logger.default.path")
	loadFiles, _ := filepath.Glob(path + "/monitor/*.log")
	loadData := passingFiles(loadFiles)
	return loadData
}

func getDirSize(path string) uint64 {
	var size int64
	wd, _ := os.Getwd()
	if !helper.IsExist(wd + path) {
		return 0
	}
	_ = filepath.Walk(wd+path, func(_ string, info os.FileInfo, err error) error {
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
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	bufferRead := bufio.NewReader(fd)
	data := make([]map[string]any, 0)
	for {
		line, err := bufferRead.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		curData := map[string]any{}
		err = json.Unmarshal([]byte(line), &curData)
		if err != nil {
			continue
		}
		data = append(data, curData)
	}
	return data, nil
}
