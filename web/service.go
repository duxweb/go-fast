package web

import (
	"context"

	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/hook"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/urfave/cli/v3"
)

func Service() *service.Config {
	return &service.Config{
		Name: "web",
		Init: Init,
		Boot: Boot,
		Cmd:  Command,
	}
}

func Command() []*cli.Command {
	cmd := &cli.Command{
		Category: "service",
		Name:     "web",
		Usage:    "starting the web service",
		Action: func(context.Context, *cli.Command) error {
			// 启动服务
			global.Service.Boot()
			// 执行启动钩子
			hook.RunBoot()
			// 启动 Web 服务
			Start()
			return nil
		},
	}

	return []*cli.Command{
		cmd,
	}
}
