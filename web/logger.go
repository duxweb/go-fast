package web

import (
	"context"
	"fmt"
	"github.com/duxweb/go-fast/logger"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"io"
	"log/slog"
)

type FiberLogger struct {
	Logger    *slog.Logger
	SlogLevel slog.Level
	Ctx       context.Context
}

func (l *FiberLogger) SetLevel(level fiberlog.Level) {

}

func (l *FiberLogger) SetOutput(writer io.Writer) {
}

func (l *FiberLogger) WithContext(ctx context.Context) fiberlog.CommonLogger {
	return &FiberLogger{Logger: logger.Log(), Ctx: ctx}
}

func (l *FiberLogger) Trace(v ...interface{}) {
	l.Logger.DebugContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Debug(v ...interface{}) {
	l.Logger.DebugContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Info(v ...interface{}) {
	l.Logger.InfoContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Warn(v ...interface{}) {
	l.Logger.WarnContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Error(v ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Fatal(v ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Panic(v ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, fmt.Sprint(v...))
}

func (l *FiberLogger) Tracef(format string, v ...interface{}) {
	l.Logger.DebugContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Debugf(format string, v ...interface{}) {
	l.Logger.DebugContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Infof(format string, v ...interface{}) {
	l.Logger.InfoContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Warnf(format string, v ...interface{}) {
	l.Logger.WarnContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Errorf(format string, v ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Panicf(format string, v ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, fmt.Sprintf(format, v...))
}

func (l *FiberLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.Logger.DebugContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func (l *FiberLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Logger.DebugContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func (l *FiberLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.Logger.InfoContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func (l *FiberLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.Logger.WarnContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func (l *FiberLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func (l *FiberLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func (l *FiberLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.Logger.ErrorContext(l.Ctx, msg, slog.Any("data", keysAndValues))
}

func LoggerAdaptor() fiberlog.AllLogger {
	return &FiberLogger{Logger: logger.Log(), Ctx: context.Background()}
}
