package global

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

var (
	App          *fiber.App
	Version      = "v2.0.0"
	BootTime     time.Time
	TablePrefix  = "app_"
	Debug        bool
	DebugMsg     string
	TimeLocation = time.UTC
	Ctx          = context.Background()
	DirList      []string
	ConfigDir    = "./config/"
)
