package queue

import (
	"log/slog"
	"time"

	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-queue"
	"github.com/duxweb/go-queue/drivers/memory"
	"github.com/duxweb/go-queue/drivers/redis"
	"github.com/duxweb/go-queue/drivers/sqlite"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"
)

func RegisterConfigQueues() {
	var drivers []*koanf.Koanf
	var workers []*koanf.Koanf

	if !config.IsLoad("queue") {
		// 如果没有配置文件，手动构建默认配置
		defaultDriver := koanf.New(".")
		defaultDriver.Load(confmap.Provider(map[string]interface{}{
			"name": "default",
			"type": "memory",
		}, "."), nil)
		drivers = []*koanf.Koanf{defaultDriver}

		defaultWorker := koanf.New(".")
		defaultWorker.Load(confmap.Provider(map[string]interface{}{
			"name":        "default",
			"device_name": "default",
			"num":         20,
			"interval":    1,
			"retry":       3,
			"retry_delay": 30,
			"timeout":     300,
		}, "."), nil)
		workers = []*koanf.Koanf{defaultWorker}
	} else {
		drivers = config.Load("queue").Slices("drivers")
		workers = config.Load("queue").Slices("workers")
	}

	// 统一处理驱动配置
	for _, driverCfg := range drivers {
		name := driverCfg.String("name")
		name = lo.Ternary(name == "", "default", name)
		Type := driverCfg.String("type")
		Type = lo.Ternary(Type == "", "memory", Type)

		var driver queue.QueueDriver
		var err error
		switch Type {
		case "memory":
			driver = memory.New()
		case "sqlite":
			driver, err = sqlite.New(&sqlite.SQLiteOptions{
				// 数据库文件路径
				DBPath: driverCfg.String("path"),
			})
			if err != nil {
				slog.Error("Failed to queue SQLite driver", "error", err)
			}
		case "redis":
			// 使用 Redis 数据库
			client := database.Redis(driverCfg.String("database"))
			driver, err = redis.New(redis.WithClient(client))
		default:
			slog.Error("Invalid queue type", "type", Type)
		}
		if err != nil {
			slog.Error("Failed to register queue driver", "error", err)
			return
		}

		Queue().RegisterDriver(name, driver)
	}

	// 统一处理工作器配置
	for _, worker := range workers {
		workerConfig := &queue.WorkerConfig{
			DeviceName: worker.String("device_name"),
			Num:        worker.Int("num"),                            // Concurrent workers / 并发工作数量
			Interval:   worker.Duration("interval") * time.Second,    // Polling interval / 轮询间隔
			Retry:      worker.Int("retry"),                          // Retry attempts / 重试次数
			RetryDelay: worker.Duration("retry_delay") * time.Second, // Delay between retries / 重试间隔
			Timeout:    worker.Duration("timeout") * time.Second,     // Task timeout / 任务超时时间
		}

		if workerConfig.DeviceName == "" {
			slog.Error("Queue worker name cannot be empty")
			continue
		}

		name := worker.String("name")
		if name == "" {
			slog.Error("Queue worker name is required")
			continue
		}
		Queue().RegisterWorker(worker.String("name"), workerConfig)
	}
}
