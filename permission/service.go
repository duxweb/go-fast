package permission

import (
	"context"
	"os"

	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/hook"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v3"
)

func Service() *service.Config {
	return &service.Config{
		Name: "permission",
		Cmd:  Command,
	}
}

func Command() []*cli.Command {
	cmd := &cli.Command{
		Category: "permission",
		Name:     "permission:list",
		Usage:    "View all permissions registered",
		Action: func(context.Context, *cli.Command) error {
			// 启动服务
			global.Service.Register()
			// 执行启动钩子
			hook.RunBoot()

			for name, list := range Permissions {
				color.Println(name)
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Name"})
				for _, item := range list.Get() {
					t.AppendRow(table.Row{item["name"]})
					t.AppendSeparator()

					rows := make([]table.Row, 0)
					if children, ok := item["children"]; ok {
						for _, m := range children.([]map[string]any) {
							rows = append(rows, table.Row{m["name"]})
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
