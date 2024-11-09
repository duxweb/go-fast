package web

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/route"
	"github.com/duxweb/go-fast/task"
	"github.com/duxweb/go-fast/task/queue"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
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

			// 启动队列服务
			task.Queue().Start()
			task.Cron().Start()

			// 启动 web 服务
			Start()

			task.Queue().Add("default", queue.QueueAdd{
				Name: "ping",
			})

			<-ctx.Done()

			ctx, cancel := context.WithTimeout(global.CtxBackground, 3*time.Second)
			defer cancel()

			if err := global.App.Shutdown(ctx); err != nil {
				color.Errorln(err.Error())
			}

			if err := global.Injector.Shutdown(); err != nil {
				color.Errorln("Stop service")
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
				t.AppendHeader(table.Row{"Name", "Method", "Path", "Label"})
				rows := make([]table.Row, 0)

				for _, item := range list.ParseData(nil, list.Prefix) {
					rows = append(rows, table.Row{item["name"], item["method"], item["path"], item["label"]})
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
