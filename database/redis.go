package database

import (
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
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
	dbConfig := config.Load("database").GetStringMap("redis.drivers")
	for name, _ := range dbConfig {
		do.ProvideNamed[*RedisService](global.Injector, "redis."+name, func(injector do.Injector) (*RedisService, error) {
			return NewRedis(name), nil
		})
	}

	clusterConfig := config.Load("database").GetStringMap("redisCluster.drivers")
	for name, _ := range clusterConfig {
		do.ProvideNamed[*RedisClusterService](global.Injector, "redisCluster."+name, func(injector do.Injector) (*RedisClusterService, error) {
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
	err := config.Load("database").ReadInConfig()
	if err != nil {
		return nil, err
	}
	dbConfig := config.Load("database").GetStringMapString("redis.drivers." + name)
	client := redis.NewClient(&redis.Options{
		Addr:        dbConfig["host"] + ":" + dbConfig["port"],
		Password:    dbConfig["password"],
		DB:          gocast.Number[int](dbConfig["db"]),
		DialTimeout: 5 * time.Second,
	})

	_, err = client.Ping(global.CtxBackground).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ConnCluster(name string) (*redis.ClusterClient, error) {
	err := config.Load("database").ReadInConfig()
	if err != nil {
		return nil, err
	}
	dbConfig := config.Load("database").Sub("redisCluster.drivers." + name)
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    dbConfig.GetStringSlice("addrs"),
		Username: dbConfig.GetString("username"),
		Password: dbConfig.GetString("password"),
	})
	_, err = client.Ping(global.CtxBackground).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
