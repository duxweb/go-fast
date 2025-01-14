package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logs     = map[string]*slog.Logger{}
	logMutex sync.RWMutex
)

func Log(names ...string) *slog.Logger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}

	// 先尝试读取
	logMutex.RLock()
	if t, ok := logs[name]; ok {
		logMutex.RUnlock()
		return t
	}
	logMutex.RUnlock()

	// 如果不存在，加写锁创建新的 logger
	logMutex.Lock()
	defer logMutex.Unlock()

	// 双重检查，避免在获取写锁期间其他 goroutine 已经创建了 logger
	if t, ok := logs[name]; ok {
		return t
	}

	level := config.Load("logger").GetString(name + ".level")

	parseLevel, err := log.ParseLevel(level)
	if err != nil {
		parseLevel = log.DebugLevel
	}

	logger := slog.New(
		slogmulti.Fanout(
			GetWriterHeader(
				level,
				name,
			),
			log.NewWithOptions(os.Stdout, log.Options{
				ReportCaller:    true,
				ReportTimestamp: true,
				TimeFormat:      time.DateTime,
				Prefix:          "Dux",
				Level:           parseLevel,
			}),
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
		Filename:   fmt.Sprintf(global.DataDir+"logs/%s.log", name),     // Log file path.
		MaxSize:    config.Load("logger").GetInt("default.max_size"),    // Maximum size of each log file to be saved, unit: M.
		MaxBackups: config.Load("logger").GetInt("default.max_backups"), // Number of file backups.
		MaxAge:     config.Load("logger").GetInt("default.max_age"),     // Maximum number of days to keep the files.
		Compress:   config.Load("logger").GetBool("default.compress"),   // Compression status.
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
