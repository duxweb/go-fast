package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	coreLogger "github.com/duxweb/go-fast/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var MigrateModel = make([]any, 0)

func GormMigrate(dst ...any) {
	MigrateModel = append(MigrateModel, dst...)
}

type GormService struct {
	engine *gorm.DB
}

func (s *GormService) Shutdown() error {
	sqlDB, err := s.engine.DB()
	if err != nil {
		// log
	}
	return sqlDB.Close()
}

func Gorm() *gorm.DB {
	return do.MustInvoke[*GormService](global.Injector).engine
}

func GormInit() {

	dbConfig := config.Load("database").GetStringMapString("db")

	var connect gorm.Dialector
	if dbConfig["type"] == "mysql" {
		connect = mysql.Open(dbConfig["username"] + ":" + dbConfig["password"] + "@tcp(" + dbConfig["host"] + ":" + dbConfig["port"] + ")/" + dbConfig["dbname"] + "?charset=utf8mb4&parseTime=True&loc=Local")
	}
	if dbConfig["type"] == "postgresql" {
		connect = postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbConfig["host"],
			dbConfig["username"],
			dbConfig["password"],
			dbConfig["dbname"],
			dbConfig["port"],
		))
	}
	if dbConfig["type"] == "sqlite" {
		connect = sqlite.Open(dbConfig["file"])
	}
	database, err := gorm.Open(connect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   global.TablePrefix,
			SingularTable: true,
		},
		Logger: GormLogger(),
	})
	if err != nil {
		panic("database error: " + err.Error())
	}

	do.ProvideValue[*GormService](global.Injector, &GormService{
		engine: database,
	})

	// Set Connection Pool
	sqlDB, err := do.MustInvoke[*gorm.DB](nil).DB()
	if err != nil {
		panic("database error: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.Load("app").GetInt("database.maxIdleConns"))
	sqlDB.SetMaxOpenConns(config.Load("app").GetInt("database.maxOpenConns"))

}

// Logger Database Logs
type Logger struct {
	SlowThreshold             time.Duration
	SourceField               string
	IgnoreRecordNotFoundError bool
	Logger                    zerolog.Logger
	LogLevel                  gormLogger.LogLevel
}

func GormLogger() *Logger {
	vLog := coreLogger.New(
		coreLogger.GetWriter(
			config.Load("app").GetString("logger.db.level"),
			"gorm",
			true,
		)).With().Caller().CallerWithSkipFrameCount(5).Timestamp().Logger()

	return &Logger{
		SlowThreshold:             1 * time.Second,
		Logger:                    vLog,
		LogLevel:                  gormLogger.Silent,
		IgnoreRecordNotFoundError: true,
	}
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &Logger{
		Logger:                    l.Logger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel < gormLogger.Info {
		return
	}
	l.Logger.Info().Msgf(s, args)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel < gormLogger.Warn {
		return
	}
	l.Logger.Warn().Msgf(s, args)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel < gormLogger.Error {
		return
	}
	l.Logger.Error().Msgf(s, args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := map[string]interface{}{
		"sql":      sql,
		"duration": elapsed,
	}
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.Logger.Error().Err(err).Fields(fields).Msg("[GORM] query error")
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		l.Logger.Warn().Fields(fields).Msgf("[GORM] slow query")
	case l.LogLevel >= gormLogger.Info:
		l.Logger.Debug().Fields(fields).Msgf("[GORM] query")
	}
}
