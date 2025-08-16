package database

import (
	"context"

	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/hook"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/gookit/color"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"
)

func Service() *service.Config {
	return &service.Config{
		Name: "database",
		Init: Init,
		Cmd:  Command,
	}
}

func Init() error {
	GormInit()
	RedisInit()
	MongoInit()
	return nil
}

func Command() []*cli.Command {
	cmd := &cli.Command{
		Category: "database",
		Name:     "db:sync",
		Usage:    "Synchronous database structure",
		Action: func(context.Context, *cli.Command) error {
			// 启动服务
			global.Service.Boot()
			// 执行启动钩子
			hook.RunBoot()

			err := SyncDatabase()

			if err != nil {
				color.Println(err.Error())
			}

			return nil
		},
	}

	return []*cli.Command{
		cmd,
	}
}

func SyncDatabase() error {
	models := make([]any, 0)
	sends := make([]func(db *gorm.DB), 0)

	for _, model := range MigrateModel {
		if m, ok := model.(Migrate); ok {
			hasTable := Gorm().Migrator().HasTable(m.Model)
			if !hasTable {
				sends = append(sends, m.Seed)
			}
			models = append(models, m.Model)
		} else {
			models = append(models, model)
		}
	}

	err := Gorm().AutoMigrate(models...)
	if err != nil {
		return err
	}

	for _, send := range sends {
		send(Gorm())
	}
	return nil
}
