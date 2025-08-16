package lock

import (
	"log"

	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/global"
	golock "github.com/duxweb/go-lock"
	"github.com/duxweb/go-lock/drivers"
	"github.com/samber/do/v2"
)

type LockService struct {
	manager *golock.Manager
	name    string
}

func (s *LockService) Shutdown() error {
	return nil
}

func LockInit() {
	var lockConfig []string
	if config.IsLoad("lock") {
		lockConfig = config.Load("lock").MapKeys("drivers")
	} else {
		lockConfig = []string{"default"}
	}

	for _, name := range lockConfig {
		do.ProvideNamed(global.Injector, "lock."+name, func(injector do.Injector) (*LockService, error) {
			return NewLock(name), nil
		})
	}
}

// Lock 获取指定名称的锁管理器，默认为 "default"
func Lock(name ...string) *golock.Manager {
	n := "default"
	if len(name) > 0 {
		n = name[0]
	}
	service := do.MustInvokeNamed[*LockService](global.Injector, "lock."+n)
	return service.manager
}

// NewLock 创建新的锁服务
func NewLock(name string) *LockService {

	var lockConfig map[string]string

	if config.IsLoad("lock") {
		lockConfig = config.Load("lock").StringMap("drivers." + name)
	} else {
		lockConfig = map[string]string{
			"type": "memory",
		}
	}

	var provider golock.LockProvider

	lockType := lockConfig["type"]
	if lockType == "" {
		lockType = "memory"
	}

	var err error

	switch lockType {
	case "redis":
		// 使用 Redis 数据库
		client := database.Redis(lockConfig["database"])

		provider, err = drivers.NewRedisDriver(&drivers.RedisOptions{
			Client: client,
		})
		if err != nil {
			log.Fatalf("Failed to create Redis lock provider: %v", err)
		}
	default:
		// 默认使用内存锁
		provider = drivers.NewMemoryDriver()
	}

	manager := golock.New(provider)

	return &LockService{
		manager: manager,
		name:    name,
	}
}
