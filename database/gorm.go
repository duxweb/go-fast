package database

import (
	"context"
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	coreLogger "github.com/duxweb/go-fast/logger"
	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/samber/do/v2"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Migrate struct {
	Model any
	Seed  func(db *gorm.DB)
}

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

func GormInit() {
	dbConfig := config.Load("database").GetStringMap("db.drivers")
	for name, _ := range dbConfig {
		do.ProvideNamed[*GormService](global.Injector, "orm."+name, func(injector do.Injector) (*GormService, error) {
			return NewGorm(name), nil
		})
	}
}

func Gorm(name ...string) *gorm.DB {
	n := "default"
	if len(name) > 0 {
		n = name[0]
	}
	client := do.MustInvokeNamed[*GormService](global.Injector, "orm."+n)
	return client.engine
}

func GormCtx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		return tx
	}
	return Gorm().WithContext(ctx)
}

func NewGorm(name string) *GormService {
	// 重新读取服务
	err := config.Load("database").ReadInConfig()
	if err != nil {
		panic("redis error :" + err.Error())
	}
	dbConfig := config.Load("database").GetStringMapString("db.drivers." + name)
	var connect gorm.Dialector
	if dbConfig["type"] == "mysql" {
		connect = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig["username"],
			dbConfig["password"],
			dbConfig["host"],
			dbConfig["port"],
			dbConfig["database"],
		))
	}
	if dbConfig["type"] == "postgresql" {
		connect = postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbConfig["host"],
			dbConfig["username"],
			dbConfig["password"],
			dbConfig["database"],
			dbConfig["port"],
		))
	}
	if dbConfig["type"] == "sqlite" {
		connect = sqlite.Open(dbConfig["file"] + "?_pragma=journal_mode(WAL)")
	}
	database, err := gorm.Open(connect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   global.TablePrefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   slogGorm.New(slogGorm.WithHandler(coreLogger.GetWriterHeader(config.Load("logger").GetString("db.level"), "db"))),
	})
	if err != nil {
		panic("database error: " + err.Error())
	}

	// Set Connection Pool
	sqlDB, err := database.DB()
	if err != nil {
		panic("database error: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(cast.ToInt(dbConfig["maxIdleConns"]))
	sqlDB.SetMaxOpenConns(cast.ToInt(dbConfig["maxOpenConns"]))

	return &GormService{
		engine: database,
	}
}

func SwitchGorm(name string) error {
	// 关闭原服务
	err := do.ShutdownNamed(global.Injector, "orm."+name)
	if err != nil {
		return err
	}
	// 替换服务
	do.OverrideNamed[*GormService](global.Injector, "orm."+name, func(injector do.Injector) (*GormService, error) {
		return NewGorm(name), nil
	})
	return nil
}
