package database

import (
	"time"

	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	"github.com/spf13/cast"
)

type RedisService struct {
	engine *redis.Client
}

type RedisClusterService struct {
	engine *redis.ClusterClient
}

func (s *RedisService) Shutdown() error {
	return s.engine.Close()
}

func RedisInit() {
	dbConfig := config.Load("database").MapKeys("redis.drivers")
	for _, name := range dbConfig {
		do.OverrideNamed(global.Injector, "redis."+name, func(injector do.Injector) (*RedisService, error) {
			return NewRedis(name), nil
		})
	}

	clusterConfig := config.Load("database").MapKeys("redisCluster.drivers")
	for _, name := range clusterConfig {
		do.ProvideNamed(global.Injector, "redisCluster."+name, func(injector do.Injector) (*RedisClusterService, error) {
			return NewRedisCluster(name), nil
		})
	}
}

func Redis(name ...string) *redis.Client {
	n := "default"
	if len(name) > 0 {
		n = name[0]
	}
	client := do.MustInvokeNamed[*RedisService](global.Injector, "redis."+n)
	return client.engine
}

func NewRedis(name string) *RedisService {
	client, err := ConnSingle(name)
	if err != nil {
		panic("redis error :" + err.Error())
	}
	return &RedisService{
		engine: client,
	}
}

func RedisCluster(name ...string) *redis.ClusterClient {
	n := "default"
	if len(name) > 0 {
		n = name[0]
	}
	client := do.MustInvokeNamed[*RedisClusterService](global.Injector, "redisCluster."+n)
	return client.engine
}

func NewRedisCluster(name string) *RedisClusterService {
	client, err := ConnCluster(name)
	if err != nil {
		panic("redis error :" + err.Error())
	}
	return &RedisClusterService{
		engine: client,
	}
}

func ConnSingle(name string) (*redis.Client, error) {

	var clientConfig struct {
		Host     string        `koanf:"host"`
		Port     int           `koanf:"port"`
		Password string        `koanf:"password"`
		DB       int           `koanf:"db"`
		Timeout  time.Duration `koanf:"timeout"`
	}

	err := config.Load("database").Unmarshal("redis.drivers."+name, &clientConfig)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:        clientConfig.Host + ":" + cast.ToString(clientConfig.Port),
		Password:    clientConfig.Password,
		DB:          clientConfig.DB,
		DialTimeout: clientConfig.Timeout * time.Second,
	})

	_, err = client.Ping(global.Ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ConnCluster(name string) (*redis.ClusterClient, error) {

	var clientConfig struct {
		Addrs    []string `koanf:"addrs"`
		Username string   `koanf:"username"`
		Password string   `koanf:"password"`
		DB       int      `koanf:"db"`
	}

	err := config.Load("database").Unmarshal("redisCluster.drivers."+name, &clientConfig)
	if err != nil {
		return nil, err
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    clientConfig.Addrs,
		Username: clientConfig.Username,
		Password: clientConfig.Password,
	})
	_, err = client.Ping(global.Ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
