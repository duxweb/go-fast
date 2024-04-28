package service

import (
	"github.com/duxweb/go-fast/cache"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/validator"
	"github.com/duxweb/go-fast/views"
	"github.com/samber/do"
)

func Init() {
	global.Injector = do.New()
	config.Init()
	logger.Init()
	cache.Init()
	i18n.Init()
	validator.Init()
	views.Init()
}
