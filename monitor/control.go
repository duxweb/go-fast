package monitor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/duxweb/go-fast/global"
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
