package database

import (
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"github.com/spf13/cast"
)

func RedisInit() {
	dbConfig := config.Get("database").GetStringMapString("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       cast.ToInt(dbConfig["db"]),
	})
	_, err := client.Ping(global.Ctx).Result()
	if err != nil {
		panic(err.Error())
	}
	do.ProvideValue[*redis.Client](nil, client)
}
