package lock

import (
	"time"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/lock/driver"
	"github.com/samber/do/v2"
)

type LockDriver interface {
	Create(key string, ttl time.Duration) driver.LockDriver
}

func Get(name ...string) LockDriver {
	lockDriver := config.Load("use").GetString("lock.driver")
	if len(name) > 0 {
		lockDriver = name[0]
	}
	lockDriverName := "lock." + lockDriver

	if _, err := do.InvokeNamed[LockDriver](global.Injector, lockDriverName); err != nil {
		do.ProvideNamed[LockDriver](global.Injector, lockDriverName, func(i do.Injector) (LockDriver, error) {
			var s LockDriver
			switch lockDriver {
			case "redis":
				s = driver.NewRedisDriver(database.Redis())
			case "memory":
			default:
				s = driver.NewMemoryDriver()
			}
			return s, nil
		})
	}

	return do.MustInvokeNamed[LockDriver](global.Injector, lockDriverName)
}

func New(driver LockDriver) *Manager {
	return &Manager{
		driver: driver,
	}
}

type Manager struct {
	driver LockDriver
}

func (m *Manager) Create(key string, ttl time.Duration) driver.LockDriver {
	return m.driver.Create(key, ttl)
}
