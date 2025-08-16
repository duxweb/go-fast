package config

import (
	"github.com/duxweb/go-fast/v2/service"
)

func Service() *service.Config {
	return &service.Config{
		Name: "config",
		Init: Init,
	}
}
