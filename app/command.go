package app

import (
	"github.com/duxweb/go-fast/global"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"os"
)

func Command() []*cli.Command {
	version := &cli.Command{
		Name:     "version",
		Category: "dev",
		Usage:    "View the version number",
		Action: func(cCtx *cli.Context) error {
			color.Redf("⇨ <red>%s</>", global.Version)
			return nil
		},
	}

	appList := &cli.Command{
		Name:     "app:list",
		Category: "app",
		Usage:    "viewing the application list",
		Action: func(cCtx *cli.Context) error {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Name"})
			rows := make([]table.Row, 0)
			for _, config := range List {
				rows = append(rows, table.Row{config.Name})
			}
			t.AppendRows(rows)
			t.Render()
			return nil
		},
	}

	return []*cli.Command{
		version,
		appList,
	}
}
