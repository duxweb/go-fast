package task

import (
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/samber/do/v2"
)

func Service() *service.Config {
	return &service.Config{
		Name: "cron",
		Init: Init,
	}
}

func Cron() *CronService {
	client := do.MustInvokeNamed[*CronService](global.Injector, "cron")
	return client
}

func Init() error {
	do.ProvideNamed(global.Injector, "cron", NewCron)
	return nil
}

func Boot() error {
	RegisterCron()
	return nil
}
