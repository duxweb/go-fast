package cache

import (
	"github.com/duxweb/go-fast/v2/service"
)

func Service() *service.Config {
	return &service.Config{
		Name: "cache",
		Init: func() error {
			CacheInit()
			return nil
		},
	}
}
