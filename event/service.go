package event

import (
	"context"
	"os"
	"sort"

	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/hook"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v3"
)

func Service() *service.Config {
	return &service.Config{
		Name: "event",
		Init: Init,
		Cmd:  Command,
	}
}

func Init() error {
	Register()
	return nil
}

func Command() []*cli.Command {
	cmd := &cli.Command{
		Category: "event",
		Name:     "event:list",
		Usage:    "View all event listeners",
		Action: func(context.Context, *cli.Command) error {
			// 启动服务
			global.Service.Register()
			// 执行启动钩子
			hook.RunBoot()

			// 如果没有事件监听器
			if len(Listener) == 0 {
				color.Yellow.Println("⚠ No event listeners registered")
				return nil
			}

			// 创建表格
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Event Name", "Listener Count", "Locations"})

			// 获取所有事件名并排序
			eventNames := make([]string, 0, len(Listener))
			for name := range Listener {
				eventNames = append(eventNames, name)
			}
			sort.Strings(eventNames)

			// 统计总数
			totalListeners := 0

			// 添加行
			rows := make([]table.Row, 0)
			for _, name := range eventNames {
				locations := Listener[name]
				totalListeners += len(locations)

				// 将多个位置用换行符连接
				locationStr := ""
				for i, loc := range locations {
					if i > 0 {
						locationStr += "\n"
					}
					locationStr += loc
				}

				rows = append(rows, table.Row{
					name,
					len(locations),
					locationStr,
				})
			}

			t.AppendRows(rows)

			// 添加统计行
			t.AppendSeparator()
			t.AppendRow(table.Row{
				color.Bold.Sprint("Total"),
				color.Bold.Sprint(totalListeners),
				color.Bold.Sprintf("%d events", len(eventNames)),
			})

			t.Render()

			return nil
		},
	}

	return []*cli.Command{
		cmd,
	}
}
