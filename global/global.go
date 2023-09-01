package global

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"golang.org/x/text/language"
	"time"
)

var (
	App         *fiber.App
	Version     = "v2.0.0"
	BootTime    time.Time
	TablePrefix = "app_"
	Lang        = language.English

	Injector *do.Injector

	Debug        bool
	DebugMsg     string
	Ctx          context.Context
	TimeLocation = time.UTC
	DirList      []string
	ConfigDir    = "./config/"
)
