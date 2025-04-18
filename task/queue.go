package task

import (
	"context"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/task/queue"
	"github.com/gookit/color"
	"github.com/samber/do/v2"
)

type QueueService struct {
	Queue queue.Queue
}

func (s *QueueService) Shutdown() error {
	if s == nil {
		return nil
	}
	return s.Queue.Close()
}

func QueueInit() {
	do.ProvideNamed(global.Injector, "queue", NewQueue)
}

func NewQueue(i do.Injector) (*QueueService, error) {
	dbConfig := config.Load("use").GetString("queue.driver")
	if dbConfig == "" {
		dbConfig = "base"
	}

	var driver queue.Queue

	switch dbConfig {
	case "base":
		driver = queue.NewBase()
	case "redis":
		driver = queue.NewRedis()
	}

	driver.Worker("default")

	workers := config.Load("use").GetStringSlice("queue.workers")
	for _, worker := range workers {
		driver.Worker(worker)
	}

	driver.Register("default", "ping", func(ctx context.Context, params []byte) error {
		color.Println("⇨ <green>Queue service ping</>")
		return nil
	})

	return &QueueService{
		Queue: driver,
	}, nil
}

func Queue() queue.Queue {
	client := do.MustInvokeNamed[*QueueService](global.Injector, "queue")
	return client.Queue
}
