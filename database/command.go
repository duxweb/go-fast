package database

import (
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func Command() []*cli.Command {
	sync := &cli.Command{
		Name:     "db:sync",
		Usage:    "Synchronous database structure",
		Category: "database",
		Action: func(ctx *cli.Context) error {
			Register()

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
				color.Println(err.Error())
				return err
			}

			for _, send := range sends {
				send(Gorm())
			}

			return nil
		},
	}
	return []*cli.Command{
		sync,
	}
}
