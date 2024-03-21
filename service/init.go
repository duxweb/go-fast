package service

import (
	"context"
	"github.com/duxweb/go-fast/cache"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/validator"
	"github.com/duxweb/go-fast/views"
	"github.com/samber/do"
)

var Server = ServerStatus{}
var ContextCancel context.CancelFunc

type ServerStatus struct {
	Database bool
	Redis    bool
	Mongodb  bool
}

func Init() {
	global.Ctx, ContextCancel = context.WithCancel(context.Background())
	global.Injector = do.New()
	config.Init()
	logger.Init()
	cache.Init()
	i18n.Init()
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
