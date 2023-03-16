package service

import (
	"github.com/duxweb/go-fast/cache"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/validator"
	"github.com/duxweb/go-fast/views"
)

var Server = ServerStatus{}

type ServerStatus struct {
	Database bool
	Redis    bool
	Mongodb  bool
}

func Init() {
	config.Init()
	logger.Init()
	cache.Init()
	validator.Init()
	views.Init()
	if Server.Database {
		database.GormInit()
	}
	if Server.Redis {
		database.RedisInit()
	}
	if Server.Mongodb {
		database.QmgoInit()
	}

}
