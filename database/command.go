package database

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func Command(command *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "db:sync",
		Short: "Synchronous database structure",
		Run: func(cmd *cobra.Command, args []string) {
			for _, model := range MigrateModel {
				err := Gorm().AutoMigrate(model)
				if err != nil {
					color.Println(err.Error())
				}
			}
		},
	}
	command.AddCommand(cmd)
}
