package service

import "github.com/urfave/cli/v3"

type Config struct {
	Name     string
	Init     func() error
	Register func() error
	Shutdown func() error
	Cmd      func() []*cli.Command
}
