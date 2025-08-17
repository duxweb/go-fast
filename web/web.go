package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/logger"
	"github.com/gookit/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func Init() error {
	domain := "http://0.0.0.0:8900"
	host := "0.0.0.0"
	port := "8900"
	if config.IsLoad("use") {
		if config.Load("use").Exists("server.port") {
			port = config.Load("use").String("server.port")
		}
		if config.Load("use").Exists("server.host") {
			host = config.Load("use").String("server.host")
		}
		domain = fmt.Sprintf("http://%s:%s", host, port)
	}

	global.Web = global.WebConfig{
		Port:   port,
		Host:   host,
		Domain: domain,
	}

	router := echo.New()
	router.Debug = global.Debug
	router.HideBanner = true
	router.HidePort = true

	router.Renderer = ViewHandler()
	router.Logger = LoggerHandler()
	router.HTTPErrorHandler = ErrorHandler()
	router.IPExtractor = IpHandler()

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	router.Use(middleware.Recover())
	router.Use(middleware.RequestID())
	router.Use(slogecho.NewWithFilters(logger.Log("web"), slogecho.Accept(func(c echo.Context) bool {
		if strings.HasPrefix(c.Path(), "/static") || strings.HasPrefix(c.Path(), "/docs") || strings.HasPrefix(c.Path(), "/openapi.yaml") {
			return false
		}
		return true
	})))
	router.Use(I18n())
	router.Use(middleware.Gzip())

	global.Router = router

	router.GET("/ws", WebsocketHandler())

	return nil
}

func Register() error {
	global.Router.Group("/").Static("", "./public")

	if global.StaticFs != nil {
		global.Router.StaticFS("/static", echo.MustSubFS(*global.StaticFs, "static"))
	}

	global.Router.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "template/welcome.html", nil)
	})

	return nil
}

func Start() {

	global.GetBanner()
	global.BootTime = time.Now()
	color.Println(fmt.Sprintf("⇨ <green>Server start http://%s:%s</>\n", global.Web.Host, global.Web.Port))

	err := http.ListenAndServe(fmt.Sprintf(":%s", global.Web.Port), global.Router)
	if err != nil {
		color.Errorln("server start error:", err.Error())
	}

}
