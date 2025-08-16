package cmd

import "github.com/urfave/cli/v3"

// 服务配置
// Service configuration
var (
	Commands = make([]*cli.Command, 0)
)

type Command struct {
}

// New 创建命令
// New create command
func New() *Command {
	return &Command{}
}

func Register(opt []*cli.Command) {
	Commands = append(Commands, opt...)
}
