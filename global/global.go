package global

import (
	"context"
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"time"
)

var (
	App         *echo.Echo
	Version     = "v0.0.1"
	BootTime    time.Time
	TablePrefix = "app_"
	Lang        = "en-US"

	Injector *do.Injector

	Debug         bool
	DebugMsg      string
	CtxBackground = context.Background()
	TimeLocation  = time.UTC
	DirList       []string
	ConfigDir     = "./config/"

	StaticFs *embed.FS
)
