package resources

import (
	"github.com/duxweb/go-fast/v2/service"
)

func Service() *service.Config {
	return &service.Config{
		Name:     "resources",
		Register: RegisterService,
	}
}

func RegisterService() error {
	Register()
	return nil
}
