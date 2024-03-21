package web

import (
	"context"
	"github.com/duxweb/go-fast/app"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/monitor"
	"github.com/duxweb/go-fast/service"
	"github.com/duxweb/go-fast/task"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Command() []*cli.Command {
	cmd := &cli.Command{
		Name:  "web",
		Usage: "starting the web service",
		Action: func(cCtx *cli.Context) error {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch,
				os.Interrupt,
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM)

			service.Init()
			task.Init()
			Init()
			monitor.Init()
			app.Init()

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

			<-ch
			service.ContextCancel()
			err := global.Injector.Shutdown()
			if err != nil {
				color.Errorln("Stop service")
			}
			//color.Println("⇨ <orange>Stop websocket</>")
			//websocket.Release()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := global.App.Shutdown(ctx); err != nil {
				color.Errorln(err.Error())
			}

			return nil
		},
	}

	return []*cli.Command{
		cmd,
	}
}
