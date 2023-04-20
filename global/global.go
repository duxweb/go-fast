package global

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
	"time"
)

var (
	App         *fiber.App
	Version     = "v2.0.0"
	BootTime    time.Time
	TablePrefix = "app_"
	Lang        = language.English

	Debug        bool
	DebugMsg     string
	TimeLocation = time.UTC
	Ctx          = context.Background()
	DirList      []string
	ConfigDir    = "./config/"
)
