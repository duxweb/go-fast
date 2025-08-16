package system

import (
	"github.com/duxweb/go-fast/v2/app"
	"github.com/duxweb/go-fast/v2/resources"
	"github.com/duxweb/go-fast/v2/route"
)

var config = struct {
}{}

func App() *app.Config {
	return &app.Config{
		Name:     "system",
		Config:   &config,
		Init:     Init,
		Register: Register,
		Boot:     Boot,
	}
}

func Init() error {
	route.SetRouter("web", route.New("web", ""))

	resources.Set("admin", resources.New("admin", "/admin"))

	return nil
}

func Register() error {
	return nil
}

func Boot() error {
	return nil
}
