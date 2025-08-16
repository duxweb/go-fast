package lock

import (
	"github.com/duxweb/go-fast/v2/service"
)

func Service() *service.Config {
	return &service.Config{
		Name: "lock",
		Init: func() error {
			LockInit()
			return nil
		},
	}
}
