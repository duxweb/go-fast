package core

import (
	"context"

	"github.com/duxweb/go-fast/v2/global"
	"github.com/gookit/color"
	"github.com/urfave/cli/v3"
)

func CoreCommand() []*cli.Command {
	version := &cli.Command{
		Name:     "version",
		Category: "dev",
		Usage:    "View the version number",
		Action: func(context.Context, *cli.Command) error {
			color.Red.Println(global.Version)
			return nil
		},
	}

	return []*cli.Command{
		version,
	}
}
