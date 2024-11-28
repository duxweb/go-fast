package driver

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	client *redis.Client
}

func NewRedisDriver(client *redis.Client) *RedisDriver {
	return &RedisDriver{
		client: client,
	}
}

func (r *RedisDriver) Create(key string, ttl time.Duration) LockDriver {
	return &RedisLockDriver{
		key:    key,
		ttl:    ttl,
		client: r.client,
	}
}

type RedisLockDriver struct {
	key        string
	ttl        time.Duration
	client     *redis.Client
	lockValues sync.Map
}

func (r *RedisLockDriver) Acquire(wait bool) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 生成唯一的锁标识
	lockValue := uuid.New().String()

	deadline := time.Now().Add(30 * time.Second)
	for {
		success := r.client.SetNX(ctx, r.key, lockValue, r.ttl).Val()
		if err := r.client.SetNX(ctx, r.key, lockValue, r.ttl).Err(); err != nil {
			return false
		}
		if success {
			r.lockValues.Store(r.key, lockValue)
			return true
		}
		if !wait {
			return false
		}

		if time.Now().After(deadline) {
			return false
		}

		// 使用指数退避算法
		backoff := time.Duration(rand.Intn(100)) * time.Millisecond
		select {
		case <-ctx.Done():
			return false
		case <-time.After(backoff):
			continue
		}
	}
}

func (r *RedisLockDriver) Release() error {
	r.lockValues.Delete(r.key)
	return nil
}
