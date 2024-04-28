package web

import (
	"fmt"
	"github.com/duxweb/go-fast/logger"
	"github.com/labstack/gommon/log"
	"io"
	"log/slog"
	"os"
)

type EchoLogger struct {
	Logger    *slog.Logger
	SlogLevel slog.Level
}

func (l *EchoLogger) Output() io.Writer {
	return os.Stdout
}

// SetOutput 设置日志记录器的输出
func (l *EchoLogger) SetOutput(w io.Writer) {
}

func (l *EchoLogger) Prefix() string {
	return ""
}

func (l *EchoLogger) SetPrefix(p string) {
}

func (l *EchoLogger) Level() log.Lvl {
	switch l.SlogLevel {
	case slog.LevelDebug:
		return log.DEBUG
	case slog.LevelInfo:
		return log.INFO
	case slog.LevelWarn:
		return log.WARN
	case slog.LevelError:
		return log.ERROR
	default:
		return log.INFO
	}
}

func (l *EchoLogger) SetHeader(h string) {
}

func (l *EchoLogger) SetLevel(v log.Lvl) {
	switch v {
	case log.DEBUG:
		l.SlogLevel = slog.LevelDebug
	case log.INFO:
		l.SlogLevel = slog.LevelInfo
	case log.WARN:
		l.SlogLevel = slog.LevelWarn
	case log.ERROR:
		l.SlogLevel = slog.LevelError
	default:
		l.SlogLevel = slog.LevelInfo
	}
}

func (l *EchoLogger) Print(i ...interface{}) {
	l.Logger.Info(fmt.Sprint(i...))
}

func (l *EchoLogger) Printf(format string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Printj(j log.JSON) {
	l.Logger.Info("", slog.Any("data", j))
}

func (l *EchoLogger) Debug(i ...interface{}) {
	l.Logger.Debug(fmt.Sprint(i...))
}

func (l *EchoLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Debugj(j log.JSON) {
	l.Logger.Debug("", slog.Any("data", j))
}

func (l *EchoLogger) Info(i ...interface{}) {
	l.Logger.Info(fmt.Sprint(i...))
}

func (l *EchoLogger) Infof(format string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Infoj(j log.JSON) {
	l.Logger.Debug("", slog.Any("data", j))
}

func (l *EchoLogger) Warn(i ...interface{}) {
	l.Logger.Warn(fmt.Sprint(i...))
}

func (l *EchoLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Warnj(j log.JSON) {
	l.Logger.Warn("", slog.Any("data", j))
}

func (l *EchoLogger) Error(i ...interface{}) {
	l.Logger.Error(fmt.Sprint(i...))
}

func (l *EchoLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Errorj(j log.JSON) {
	l.Logger.Error("", slog.Any("data", j))
}

func (l *EchoLogger) Fatal(i ...interface{}) {
	l.Logger.Error(fmt.Sprint(i...))
}

func (l *EchoLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Fatalj(j log.JSON) {
	l.Logger.Error("", slog.Any("data", j))
}

func (l *EchoLogger) Panic(i ...interface{}) {
	l.Logger.Error(fmt.Sprint(i...))
}

func (l *EchoLogger) Panicj(j log.JSON) {
	l.Logger.Error("", slog.Any("data", j))
}

func (l *EchoLogger) Panicf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func EchoLoggerHeadAdaptor() *EchoLogger {
	return &EchoLogger{Logger: logger.Log()}
}
