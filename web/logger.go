package web

import (
	"fmt"
	"github.com/duxweb/go-fast/logger"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"io"
	"os"
)

type EchoLogger struct {
	Logger *zerolog.Logger
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
	switch l.Logger.GetLevel() {
	case zerolog.DebugLevel:
		return log.DEBUG
	case zerolog.InfoLevel:
		return log.INFO
	case zerolog.WarnLevel:
		return log.WARN
	case zerolog.ErrorLevel:
		return log.ERROR
	case zerolog.FatalLevel:
		return log.ERROR
	case zerolog.PanicLevel:
		return log.ERROR
	case zerolog.NoLevel:
		return log.INFO
	case zerolog.Disabled:
		return log.OFF
	default:
		return log.INFO
	}
}

func (l *EchoLogger) SetHeader(h string) {
}

func (l *EchoLogger) SetLevel(v log.Lvl) {
	switch v {
	case log.DEBUG:
		l.Logger.Level(zerolog.DebugLevel)
	case log.INFO:
		l.Logger.Level(zerolog.InfoLevel)
	case log.WARN:
		l.Logger.Level(zerolog.WarnLevel)
	case log.ERROR:
		l.Logger.Level(zerolog.ErrorLevel)
	case log.OFF:
		l.Logger.Level(zerolog.Disabled)
	default:
		l.Logger.Level(zerolog.InfoLevel)
	}
}

func (l *EchoLogger) Print(i ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Printf(format string, args ...interface{}) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *EchoLogger) Printj(j log.JSON) {
	l.Logger.Info().Fields(j).Msg("")
}

func (l *EchoLogger) Debug(i ...interface{}) {
	l.Logger.Debug().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debug().Msgf(format, args...)
}

func (l *EchoLogger) Debugj(j log.JSON) {
	l.Logger.Debug().Fields(j).Msg("")
}

func (l *EchoLogger) Info(i ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Infof(format string, args ...interface{}) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *EchoLogger) Infoj(j log.JSON) {
	l.Logger.Info().Fields(j).Msg("")
}

func (l *EchoLogger) Warn(i ...interface{}) {
	l.Logger.Warn().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Warn().Msgf(format, args...)
}

func (l *EchoLogger) Warnj(j log.JSON) {
	l.Logger.Warn().Fields(j).Msg("")
}

func (l *EchoLogger) Error(i ...interface{}) {
	l.Logger.Error().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Error().Msgf(format, args...)
}

func (l *EchoLogger) Errorj(j log.JSON) {
	l.Logger.Error().Fields(j).Msg("")
}

func (l *EchoLogger) Fatal(i ...interface{}) {
	l.Logger.Fatal().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Fatalj(j log.JSON) {
	l.Logger.Fatal().Fields(j).Msg("")
}

func (l *EchoLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatal().Msgf(format, args...)
}

func (l *EchoLogger) Panic(i ...interface{}) {
	l.Logger.Panic().Msg(fmt.Sprint(i...))
}

func (l *EchoLogger) Panicj(j log.JSON) {
	l.Logger.Panic().Fields(j).Msg("")
}

func (l *EchoLogger) Panicf(format string, args ...interface{}) {
	l.Logger.Panic().Msgf(format, args...)
}

func EchoLoggerHeadAdaptor() *EchoLogger {
	return &EchoLogger{Logger: logger.Log("")}
}
