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

func Redis() *redis.Client {
	return do.MustInvoke[*RedisService](global.Injector).engine
}

func RedisInit() {
	dbConfig := config.Load("database").GetStringMapString("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       gocast.Number[int](dbConfig["db"]),
	})
	_, err := client.Ping(global.Ctx).Result()
	if err != nil {
		panic(err.Error())
	}
	do.ProvideValue[*RedisService](nil, &RedisService{
		engine: client,
	})
}
