package driver

import (
	"sync"
	"time"
)

// MemoryDriver 内存锁驱动
type MemoryDriver struct {
}

func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{}
}

func (r *MemoryDriver) Create(key string, ttl time.Duration) LockDriver {
	return &MemoryLockDriver{
		locks: sync.Map{},
		key:   key,
		ttl:   ttl,
	}
}

type MemoryLockDriver struct {
	locks sync.Map
	key   string
	ttl   time.Duration
}

type memLock struct {
	mutex   sync.Mutex
	locked  bool
	timeout time.Time
}

func (r *MemoryLockDriver) Acquire(wait bool) bool {
	for {
		actual, _ := r.locks.LoadOrStore(r.key, &memLock{})
		lock := actual.(*memLock)

		lock.mutex.Lock()
		now := time.Now()
		if !lock.locked || now.After(lock.timeout) {
			lock.locked = true
			lock.timeout = now.Add(r.ttl)
			lock.mutex.Unlock()
			return true
		}
		lock.mutex.Unlock()

		if !wait {
			return false
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (r *MemoryLockDriver) Release() error {
	actual, ok := r.locks.Load(r.key)
	if !ok {
		return nil
	}
	lock := actual.(*memLock)
	lock.mutex.Lock()
	defer lock.mutex.Unlock()
	if !lock.locked {
		return nil
	}
	lock.locked = false
	return nil
}
