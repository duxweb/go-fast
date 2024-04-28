package monitor

import (
	"context"
	"fmt"
	"github.com/duxweb/go-fast/global"
	"github.com/hibiken/asynq"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
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
		Compress:   true,
	}
	log := slog.NewJSONHandler(r, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(log)
}

// Control monitor task
func Control(ctx context.Context, t *asynq.Task) error {
	data := GetMonitorData()
	logger.Debug("Monitor",
		slog.Float64("CpuPercent", data.CpuPercent),
		slog.Float64("CpuPercent", data.MemPercent),
		slog.Int("ThreadCount", data.ThreadCount),
		slog.Int("GoroutineCount", data.GoroutineCount),
		slog.Int64("Timestamp", data.Timestamp),
	)
	return nil
}
