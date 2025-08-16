package queue

import (
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/duxweb/go-queue"
	"github.com/samber/do/v2"
)

func Service() *service.Config {
	return &service.Config{
		Name: "queue",
		Init: Init,
		Boot: Boot,
	}
}

func Queue() *queue.Service {
	return do.MustInvokeNamed[*queue.Service](global.Injector, "queue")
}

func Init() error {
	do.ProvideNamed(global.Injector, "queue", func(i do.Injector) (*queue.Service, error) {
		cfg := &queue.Config{
			Context: global.Ctx,
		}
		q, err := queue.New(cfg)
		if err != nil {
			return nil, err
		}
		return q, nil
	})

	RegisterConfigQueues()
	return nil
}

func Boot() error {
	return Queue().Start()
}
