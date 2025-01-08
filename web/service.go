package web

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/gookit/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
)

func Init() {

	global.App = echo.New()
	global.App.Debug = global.Debug
	global.App.Renderer = ViewHandler()
	global.App.HideBanner = true
	global.App.HidePort = true

	global.App.Logger = LoggerHandler()
	global.App.HTTPErrorHandler = ErrorHandler()
	global.App.IPExtractor = IpHandler()

	global.App.Use(middleware.Recover())
	global.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))
	global.App.Use(middleware.RequestID())
	global.App.Use(RequestHandler())
	global.App.Use(I18nHandler())

	// 超时处理
	// timeout := 300 * time.Second
	// if config.IsLoad("use") {
	// 	isTimeout := config.Load("use").IsSet("server.timeout")
	// 	if isTimeout {
	// 		timeout = config.Load("use").GetDuration("server.timeout") * time.Second
	// 	}
	// }

	// global.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	// 	Timeout: timeout,
	// }))

	// 注册公共目录
	global.App.Group("/", CacheHandler("public, max-age=86400")).Static("", "./public")

	// 注册嵌入目录
	if global.StaticFs != nil {
		global.App.Group("/static", CacheHandler("static, max-age=86400")).StaticFS("/", echo.MustSubFS(*global.StaticFs, "static"))
	}

	global.App.GET("/", func(c echo.Context) error {
		c.Set("tpl", "app")
		return c.Render(http.StatusOK, "template/welcome.html", nil)
	})

	global.App.GET("/ws", WebsocketHandler())
}

func Start() {

	port := "8900"
	if config.IsLoad("use") {
		port = config.Load("use").GetString("server.port")
	}

	banner()
	global.BootTime = time.Now()
	color.Println("⇨ <green>Server start http://localhost:" + port + "</>")

	go func() {
		err := global.App.Start(":" + port)
		if err != nil {
			color.Errorln(err.Error())
		}
	}()

}

func banner() {

	var banner string
	banner += `   _____           ____ ____` + "\n"
	banner += `  / __  \__ ______/ ___/ __ \` + "\n"
	banner += ` / /_/ / /_/ /> </ (_ / /_/ /` + "\n"
	banner += `/_____/\_,__/_/\_\___/\____/  v` + global.Version + "\n"

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
		Value: lo.Ternary[string](global.Debug, "enabled", "disabled"),
	})
	sysMaps = append(sysMaps, item{
		Name:  "PID",
		Value: os.Getpid(),
	})
	sysMaps = append(sysMaps, item{
		Name:  "Routes",
		Value: len(global.App.Routes()),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	color.Println(banner)
}
