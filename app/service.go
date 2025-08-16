package app

import (
	"context"
	"os"

	"github.com/duxweb/go-fast/v2/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v3"
)

func Service() *service.Config {
	return &service.Config{
		Name: "app",

		Init: func() error {
			return nil
		},
		Boot: func() error {
			for _, name := range Indexes {
				appConfig := Apps[name]
				if appConfig.Init != nil {
					err := appConfig.Init()
					if err != nil {
						return err
					}
				}
			}

			for _, name := range Indexes {
				appConfig := Apps[name]
				if appConfig.Register != nil {
					err := appConfig.Register()
					if err != nil {
						return err
					}
				}
			}

			for _, name := range Indexes {
				appConfig := Apps[name]
				if appConfig.Boot != nil {
					err := appConfig.Boot()
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
		Cmd: func() []*cli.Command {
			appList := &cli.Command{
				Name:     "app:list",
				Category: "app",
				Usage:    "viewing the application list",
				Action: func(context.Context, *cli.Command) error {
					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.AppendHeader(table.Row{"Name"})
					rows := make([]table.Row, 0)
					for _, config := range Apps {
						rows = append(rows, table.Row{config.Name})
					}
					t.AppendRows(rows)
					t.Render()
					return nil
				},
			}

			return []*cli.Command{
				appList,
			}
		},
	}
}
