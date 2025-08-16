package route

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
		Name: "route",
		Boot: boot,
		Cmd:  Command,
	}
}

func boot() error {
	// 注册注解路由
	Register()
	return nil
}

func Command() []*cli.Command {
	cmd := &cli.Command{
		Category: "route",
		Name:     "route:list",
		Usage:    "View all routes registered",
		Action: func(context.Context, *cli.Command) error {
			// 启动服务
			global.Service.Boot()
			// 执行启动钩子
			hook.RunBoot()

			for name, list := range Routes {
				color.Println(name)
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Name", "Method", "Path"})
				rows := make([]table.Row, 0)

				for _, item := range list.ParseData(list.Prefix) {
					rows = append(rows, table.Row{item["name"], item["method"], item["path"]})
				}
				t.AppendRows(rows)
				t.Render()
			}

			return nil
		},
	}

	return []*cli.Command{
		cmd,
	}
}
