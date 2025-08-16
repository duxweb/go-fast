package src

import (
	"context"

	"github.com/duxweb/go-fast/v2/core"
)

// New 创建Dux
// New create Dux instance
func New(context ...context.Context) *core.App {
	core := core.New(context...)
	return core
}
