package web

import (
	"context"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/monitor"
	"github.com/duxweb/go-fast/route"
	"github.com/duxweb/go-fast/task"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Command() []*cli.Command {
	cmd := &cli.Command{
		Category: "service",
		Name:     "web",
		Usage:    "starting the web service",
		Action: func(cCtx *cli.Context) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt,
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM)
			defer stop()

			task.ListenerTask("dux.monitor", monitor.Control)
			task.ListenerScheduler("*/1 * * * *", "dux.monitor", map[string]any{}, task.PRIORITY_LOW)

			// 启动任务服务
			go func() {
				task.StartScheduler()
			}()
			// 启动队列服务
			go func() {
				task.Add("ping", &map[string]any{})
				task.StartQueue()
			}()

			// 启动 web 服务
			Start()

			<-ctx.Done()
			err := global.Injector.Shutdown()
			if err != nil {
				color.Errorln("Stop service")
			}
			ctx, cancel := context.WithTimeout(global.CtxBackground, 10*time.Second)
			defer cancel()
			if err = global.App.Shutdown(ctx); err != nil {
				color.Errorln(err.Error())
			}
			return nil
		},
	}

	routeList := &cli.Command{
		Name:     "route:list",
		Usage:    "viewing the route list",
		Category: "dev",
		Action: func(ctx *cli.Context) error {
			for name, list := range route.Routes {
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
		routeList,
	}
}
