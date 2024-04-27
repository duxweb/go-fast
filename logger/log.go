package logger

import (
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"time"
)

var logs = map[string]*slog.Logger{}

func Log(names ...string) *slog.Logger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	if t, ok := logs[name]; ok {
		return t
	}

	logger := slog.New(
		slogmulti.Fanout(
			tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.RFC3339,
				AddSource:  true,
				ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
					if attr.Key == "stack" {
						attr.Value = slog.AnyValue("")
					}
					return attr
				},
			}),
			GetWriterHeader(
				config.Load("app").GetString("logger.default.level"),
				name,
			),
		),
	)
	logs[name] = logger
	return logger
}

func Init() {
	slog.SetDefault(Log("default"))
}

func GetWriterHeader(level string, name string) *slog.JSONHandler {

	r := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./data/logs/%s.log", name),                // Log file path.
		MaxSize:    config.Load("app").GetInt("logger.default.maxSize"),    // Maximum size of each log file to be saved, unit: M.
		MaxBackups: config.Load("app").GetInt("logger.default.maxBackups"), // Number of file backups.
		MaxAge:     config.Load("app").GetInt("logger.default.maxAge"),     // Maximum number of days to keep the files.
		Compress:   config.Load("app").GetBool("logger.default.compress"),  // Compression status.
	}

	slogLevel := lo.Switch[string, slog.Leveler](level).
		Case("debug", slog.LevelDebug).
		Case("info", slog.LevelInfo).
		Case("warn", slog.LevelWarn).
		Case("error", slog.LevelError).
		Default(slog.LevelDebug)

	return slog.NewJSONHandler(r, &slog.HandlerOptions{
		Level: slogLevel,
	})
}
