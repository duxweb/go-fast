package database

import (
	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
)

type RedisService struct {
	engine *redis.Client
}

func (s *RedisService) Shutdown() error {
	return s.engine.Close()
}

func Redis(name ...string) *redis.Client {
	n := "default"
	if len(name) > 0 {
		n = name[0]
	}
	client, err := do.InvokeNamed[*RedisService](global.Injector, "redis."+n)
	if err != nil {
		client = NewRedis(n)
		do.ProvideNamedValue[*RedisService](global.Injector, "redis."+n, client)
	}
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
