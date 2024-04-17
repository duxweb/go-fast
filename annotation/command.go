package annotation

import (
	"github.com/urfave/cli/v2"
)

func Command() []*cli.Command {
	gen := &cli.Command{
		Name:     "annotation:gen",
		Usage:    "Generating Annotated Data",
		Category: "dev",
		Action: func(ctx *cli.Context) error {
			Run()
			return nil
		},
	}

	return []*cli.Command{
		gen,
	}
}
