package database

import (
	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
)

type RedisService struct {
	engine *redis.Client
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
	dbConfig := config.Load("database").GetStringMapString("redis.drivers." + name)
	client := redis.NewClient(&redis.Options{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       gocast.Number[int](dbConfig["db"]),
	})
	_, err := client.Ping(global.CtxBackground).Result()
	if err != nil {
		panic(err.Error())
	}
	return &RedisService{
		engine: client,
	}

}
