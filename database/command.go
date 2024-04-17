package database

import (
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

func Command() []*cli.Command {
	sync := &cli.Command{
		Name:     "db:sync",
		Usage:    "Synchronous database structure",
		Category: "database",
		Action: func(cCtx *cli.Context) error {
			for _, model := range MigrateModel {
				err := Gorm().AutoMigrate(model)
				if err != nil {
					color.Println(err.Error())
				}
			}
			return nil
		},
	}
	return []*cli.Command{
		sync,
	}
}
