package database

import (
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	coreLogger "github.com/duxweb/go-fast/logger"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/samber/do"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		Logger: slogGorm.New(slogGorm.WithHandler(coreLogger.GetWriterHeader(config.Load("app").GetString("logger.db.level"), "database"))),
	})
	if err != nil {
		panic("database error: " + err.Error())
	}

	do.ProvideValue[*GormService](global.Injector, &GormService{
		engine: database,
	})

	// Set Connection Pool
	sqlDB, err := database.DB()
	if err != nil {
		panic("database error: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.Load("app").GetInt("database.maxIdleConns"))
	sqlDB.SetMaxOpenConns(config.Load("app").GetInt("database.maxOpenConns"))

}
