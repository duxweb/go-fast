package app

import (
	"github.com/duxweb/go-fast/cache"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/monitor"
	"github.com/duxweb/go-fast/task"
	"github.com/duxweb/go-fast/validator"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/web"
	"github.com/duxweb/go-fast/websocket"
	"github.com/samber/do/v2"
)

func Start(t *Dux) {
	global.Injector = do.New()
	config.Init()
	logger.Init()
	cache.Init()
	i18n.Init()
	validator.Init()
	database.GormInit()
	database.RedisInit()
	database.MongoInit()
	views.Init()
	task.Init()
	web.Init()
	monitor.Init()
	websocket.Init()
	Init(t)
}
