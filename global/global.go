package global

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/duxweb/go-fast/v2/service"
	"github.com/gookit/color"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
)

var (
	// Router is the Chi router instance used by the Huma adapter
	Router       *echo.Echo
	Name         = "go-fast"
	Version      = "v2.0.0-alpha"
	BootTime     time.Time
	TablePrefix  = "app_"
	Lang         = "en-US"
	Injector     do.Injector
	Debug        bool
	DotEnv       = ""
	Ctx          = context.Background()
	TimeLocation = time.UTC

	ConfigDir = "./config/"
	DataDir   = "./data/"
	DirList   = []string{
		"./database",
		"./public",
		"./public/uploads",
		"./data",
		"./data/tmp",
		"./config",
		"./data/logs",
	}

	Service *service.Service

	StaticFs *embed.FS
	PageFs   *embed.FS
)

func GetBanner() {

	var banner string
	banner += `   _____           ____ ____` + "\n"
	banner += `  / __  \__ ______/ ___/ __ \` + "\n"
	banner += ` / /_/ / /_/ /> </ (_ / /_/ /` + "\n"
	banner += `/_____/\_,__/_/\_\___/\____/  ` + Version + "\n"

	type item struct {
		Name  string
		Value any
	}

	var sysMaps []item
	sysMaps = append(sysMaps, item{
		Name:  "Echo",
		Value: echo.Version,
	})
	sysMaps = append(sysMaps, item{
		Name:  "Debug",
		Value: lo.Ternary(Debug, "enabled", "disabled"),
	})
	sysMaps = append(sysMaps, item{
		Name:  "PID",
		Value: os.Getpid(),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	color.Println(banner)
}
