package global

import (
	"context"
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
	"time"
)

var (
	App           *fiber.App
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
	PageFs   *embed.FS
)
