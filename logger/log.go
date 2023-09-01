package logger

import (
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

var logs = map[string]*zerolog.Logger{}

func Log(names ...string) *zerolog.Logger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	if t, ok := logs[name]; ok {
		return t
	}
	writerList := make([]io.Writer, 0)
	writerList = append(writerList, GetWriter(
		config.Load("app").GetString("logger.default.level"),
		name,
		true,
	))
	log := New(writerList...).With().Caller().Logger()
	logs[name] = &log
	return &log
}

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func New(writers ...io.Writer) zerolog.Logger {
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	writers = append(writers, &console)
	multi := zerolog.MultiLevelWriter(writers...)
	return zerolog.New(multi)
}

func GetWriter(level string, name string, recursion bool) *LevelWriter {
	parseLevel, _ := zerolog.ParseLevel(level)
	return &LevelWriter{zerolog.MultiLevelWriter(&lumberjack.Logger{
		Filename:   fmt.Sprintf("./data/logs/%s.log", name),                // Log file path.
		MaxSize:    config.Load("app").GetInt("logger.default.maxSize"),    // Maximum size of each log file to be saved, unit: M.
		MaxBackups: config.Load("app").GetInt("logger.default.maxBackups"), // Number of file backups.
		MaxAge:     config.Load("app").GetInt("logger.default.maxAge"),     // Maximum number of days to keep the files.
		Compress:   config.Load("app").GetBool("logger.default.compress"),  // Compression status.
	}), parseLevel, recursion}
}

type LevelWriter struct {
	w         zerolog.LevelWriter
	level     zerolog.Level
	recursion bool
}

func (w *LevelWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}
func (w *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.level && w.recursion {
		return w.w.WriteLevel(level, p)
	}
	if level == w.level && !w.recursion {
		return w.w.WriteLevel(level, p)
	}
	return len(p), nil
}
