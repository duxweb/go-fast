package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init() {
	slog.SetDefault(Log("default"))
}

var (
	logs = sync.Map{}
)

func Log(names ...string) *slog.Logger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	if loggerVal, ok := logs.Load(name); ok {
		return loggerVal.(*slog.Logger)
	}

	level := "debug"
	if config.IsLoad("logger") {
		level = config.Load("logger").String(name + ".level")
	}

	logger := slog.New(
		slogmulti.Fanout(
			HandlerWriter(level, name),
			HandlerStdout(level, name),
		),
	)

	actual, _ := logs.LoadOrStore(name, logger)
	return actual.(*slog.Logger)
}

func Level(level string) slog.Level {
	return lo.Switch[string, slog.Level](level).
		Case("debug", slog.LevelDebug).
		Case("info", slog.LevelInfo).
		Case("warn", slog.LevelWarn).
		Case("error", slog.LevelError).
		Default(slog.LevelInfo)
}

func HandlerStdout(level string, name string) *log.Logger {
	parseLevel, err := log.ParseLevel(level)
	if err != nil {
		parseLevel = log.DebugLevel
	}

	return log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		CallerOffset:    2,
		TimeFormat:      time.DateTime,
		Prefix:          "Dux",
		Level:           parseLevel,
	})
}

func HandlerWriter(level string, name string) *slog.JSONHandler {

	maxSize := 10
	maxBackups := 10
	maxAge := 7
	compress := false

	if config.IsLoad("logger") {
		maxSize = config.Load("logger").Int("default.max_size")
		maxBackups = config.Load("logger").Int("default.max_backups")
		maxAge = config.Load("logger").Int("default.max_age")
		compress = config.Load("logger").Bool("default.compress")
	}

	r := &lumberjack.Logger{
		Filename:   fmt.Sprintf(global.DataDir+"logs/%s.log", name), // Log file path.
		MaxSize:    maxSize,                                         // Maximum size of each log file to be saved, unit: M.
		MaxBackups: maxBackups,                                      // Number of file backups.
		MaxAge:     maxAge,                                          // Maximum number of days to keep the files.
		Compress:   compress,                                        // Compression status.
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
