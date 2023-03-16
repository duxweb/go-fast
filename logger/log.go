package logger

import (
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func Log() *zerolog.Logger {
	return do.MustInvoke[*zerolog.Logger](nil)
}

func Init() {
	// Initialize default logs and output them according to log levels.
	writerList := make([]io.Writer, 0)
	levels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	for _, level := range levels {
		writerList = append(writerList, GetWriter(
			level,
			"default",
			level,
			false,
		))
	}
	log := New(writerList...).With().Timestamp().Caller().Logger()
	do.ProvideValue[*zerolog.Logger](nil, &log)

}

func New(writers ...io.Writer) zerolog.Logger {
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	writers = append(writers, &console)
	multi := zerolog.MultiLevelWriter(writers...)
	return zerolog.New(multi)
}

func GetWriter(level string, dirName string, name string, recursion bool) *LevelWriter {
	parseLevel, _ := zerolog.ParseLevel(level)
	return &LevelWriter{zerolog.MultiLevelWriter(&lumberjack.Logger{
		Filename:   fmt.Sprintf("./data/%s/%s.log", dirName, name),        // Log file path.
		MaxSize:    config.Get("app").GetInt("logger.default.maxSize"),    // Maximum size of each log file to be saved, unit: M.
		MaxBackups: config.Get("app").GetInt("logger.default.maxBackups"), // Number of file backups.
		MaxAge:     config.Get("app").GetInt("logger.default.maxAge"),     // Maximum number of days to keep the files.
		Compress:   config.Get("app").GetBool("logger.default.compress"),  // Compression status.
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
