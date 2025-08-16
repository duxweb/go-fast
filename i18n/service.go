package i18n

import (
	"github.com/duxweb/go-fast/v2/service"
)

func Service() *service.Config {
	return &service.Config{
		Name: "i18n",
		Init: Init,
	}
}
