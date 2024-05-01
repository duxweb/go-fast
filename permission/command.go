package permission

import (
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"os"
)

func Command() []*cli.Command {

	cmd := &cli.Command{
		Category: "dev",
		Name:     "permission:list",
		Usage:    "permission list",
		Action: func(ctx *cli.Context) error {
			for name, list := range Permissions {
				color.Println(name)
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Name", "Label"})
				for _, item := range list.Get() {
					t.AppendRow(table.Row{item["name"], item["label"]})
					t.AppendSeparator()

					rows := make([]table.Row, 0)

					if children, ok := item["children"]; ok {
						for _, m := range children.([]map[string]any) {
							rows = append(rows, table.Row{m["name"], m["label"]})
						}

					}
					t.AppendRows(rows)
					t.AppendSeparator()
				}
				t.Render()
			}

			return nil
		},
	}

	return []*cli.Command{
		cmd,
	}
}
