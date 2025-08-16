package cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/logger"
	"github.com/samber/do/v2"
)

type CacheService struct {
	cache *ristretto.Cache[string, any]
}

func (s *CacheService) Shutdown() error {
	if s.cache != nil {
		s.cache.Close()
		s.cache.Wait()
	}
	return nil
}

func CacheInit() {
	do.Provide[*CacheService](global.Injector, func(injector do.Injector) (*CacheService, error) {
		return NewCache(), nil
	})
}

func Cache() *ristretto.Cache[string, any] {
	service := do.MustInvoke[*CacheService](global.Injector)
	return service.cache
}

func NewCache() *CacheService {
	// 默认配置
	numCounters := int64(1e7) // 10 million counters = 10MB
	maxCost := int64(1 << 30) // 1GB 最大缓存大小
	bufferItems := int64(64)  // 缓冲区大小

	// 从配置文件读取
	if config.IsLoad("cache") {
		cacheConfig := config.Load("cache")

		// 读取缓存大小（单位：MB）
		if size := cacheConfig.Int("size"); size > 0 {
			maxCost = int64(size) * 1024 * 1024
		}

		// 读取计数器数量
		if counters := cacheConfig.Int64("counters"); counters > 0 {
			numCounters = counters
		}

		// 读取缓冲区大小
		if buffer := cacheConfig.Int64("buffer"); buffer > 0 {
			bufferItems = buffer
		}
	}

	// 创建 ristretto 缓存实例
	cache, err := ristretto.NewCache[string, any](&ristretto.Config[string, any]{
		NumCounters: numCounters, // 跟踪频率的键数
		MaxCost:     maxCost,     // 缓存的最大成本
		BufferItems: bufferItems, // 缓冲区大小
		OnEvict: func(item *ristretto.Item[any]) {
			// 可选：记录被驱逐的项
			if global.Debug {
				logger.Log("cache").Debug("cache item evicted", "key", item.Key)
			}
		},
	})

	if err != nil {
		panic("failed to create cache: " + err.Error())
	}

	return &CacheService{
		cache: cache,
	}
}

// Get 从缓存获取值
func Get[T any](key string) (T, bool) {
	var zero T
	value, found := Cache().Get(key)
	if !found {
		return zero, false
	}

	v, ok := value.(T)
	if !ok {
		return zero, false
	}

	return v, true
}

// Set 设置缓存值
func Set(key string, value any, ttl ...time.Duration) bool {
	cost := int64(1) // 默认成本为1

	// 如果提供了 TTL
	if len(ttl) > 0 && ttl[0] > 0 {
		return Cache().SetWithTTL(key, value, cost, ttl[0])
	}

	return Cache().Set(key, value, cost)
}

// Del 删除缓存
func Del(key string) {
	Cache().Del(key)
}

// Clear 清空所有缓存
func Clear() {
	Cache().Clear()
}

// Has 检查缓存是否存在
func Has(key string) bool {
	_, found := Cache().Get(key)
	return found
}

// GetOrSet 获取或设置缓存
func GetOrSet[T any](key string, fn func() (T, error), ttl ...time.Duration) (T, error) {
	// 尝试从缓存获取
	if value, found := Get[T](key); found {
		return value, nil
	}

	// 不存在则执行函数获取值
	value, err := fn()
	if err != nil {
		var zero T
		return zero, err
	}

	// 设置到缓存
	Set(key, value, ttl...)

	return value, nil
}
