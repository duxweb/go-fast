package global

import (
	"context"
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"time"
)

var (
	App           *echo.Echo
	Version       = "0.0.1"
	BootTime      time.Time
	TablePrefix   = "app_"
	Lang          = "en-US"
	Injector      do.Injector
	Debug         bool
	CtxBackground = context.Background()
	TimeLocation  = time.UTC
	DirList       []string
	ConfigDir     = "./config/"
	DataDir       = "./data/"

	StaticFs *embed.FS
)
