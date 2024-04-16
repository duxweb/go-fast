package web

import (
	"errors"
	"fmt"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/logger"
	"github.com/gookit/color"
	"github.com/samber/lo"
)

func Init() {
	global.App = echo.New()
	global.App.Debug = config.Load("app").GetBool("app.debug")
	global.App.Renderer = views.Render()
	global.App.HideBanner = true
	global.App.HidePort = true

	// 注册异常处理
	global.App.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var msg any
		var e *echo.HTTPError
		if errors.As(err, &e) {
			// http error
			code = e.Code
			msg = e.Message
		} else {
			// Other error
			logger.Log().Error().Str("stack", fmt.Sprintf("%+v", err)).Msg(err.Error())
			msg = lo.Ternary[string](global.Debug, i18n.Trans.Get("common.error"), err.Error())
		}

		if isAsync(c) {
			err = c.JSON(code, map[string]any{
				"code":    code,
				"message": msg,
			})
			if err != nil {
				logger.Log().Error().Err(err).Send()
			}
			return
		}

		c.Set("tpl", "app")

		if code == http.StatusNotFound {
			err = c.Render(http.StatusNotFound, "404.html", nil)
		} else {
			err = c.Render(http.StatusNotFound, "500.html", map[string]any{
				"code":    code,
				"message": msg,
			})
		}
		if err != nil {
			logger.Log().Error().Err(err).Send()
		}
	}

	// 异常恢复
	global.App.Use(middleware.Recover())

	// IP 获取规则
	global.App.IPExtractor = func(req *http.Request) string {
		remoteAddr := req.RemoteAddr
		if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
			remoteAddr = ip
		} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}
		if remoteAddr == "::1" {
			remoteAddr = "127.0.0.1"
		}
		return remoteAddr
	}

	// 跨域处理
	global.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	// 注册框架日志
	global.App.Logger = EchoLoggerHeadAdaptor()

	// 注册静态路由
	global.App.Static("/uploads", "./uploads")

	// 注册公共目录
	if global.StaticFs != nil {
		global.App.StaticFS("/", echo.MustSubFS(global.StaticFs, "public"))
	}

	// 注册请求id
	global.App.Use(middleware.RequestID())

	// 请求日志
	global.App.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogHost:      true,
		LogStatus:    true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogError:     true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Log("request").WithLevel(lo.Ternary[zerolog.Level](v.Latency > 1*time.Second, zerolog.WarnLevel, zerolog.InfoLevel)).Str("id", v.RequestID).
				Int("status", v.Status).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("ip", v.RemoteIP).
				Dur("latency", v.Latency).
				Err(v.Error).
				Msg("request")
			return nil
		},
	}))
}

func Start() {

	global.App.GET("/", func(c echo.Context) error {
		c.Set("tpl", "app")
		return c.Render(http.StatusOK, "welcome.html", nil)
	})

	// websocket
	global.App.GET("/ws", func(c echo.Context) error {
		token := c.QueryParam("token")
		app := c.QueryParam("app")
		if token == "" {
			logger.Log("websocket").Debug().Str("token", token).Msg("Token Not Found")
			return response.Send(c, response.Data{
				Message: "token does not exist",
			})
		}
		if app == "" {
			logger.Log("websocket").Debug().Str("app", token).Msg("App Not Found")
			return response.Send(c, response.Data{
				Message: "app does not exist",
			})
		}
		c.Request().Header.Set("token", cast.ToString(token))
		err := websocket.Service.Websocket.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return response.Send(c, response.Data{
				Message: err.Error(),
			})
		}
		return nil
	})

	port := config.Load("app").GetString("server.port")
	banner()
	global.BootTime = time.Now()
	color.Println("⇨ <green>Server start http://0.0.0.0:" + port + "</>")

	go func() {
		err := global.App.Start(":" + port)
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		if err != nil {
			color.Errorln(err.Error())
		}
	}()

}

func banner() {
	debugBool := config.Load("app").GetBool("server.debug")

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
		Value: lo.Ternary[string](debugBool, "enabled", "disabled"),
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

func isAsync(ctx echo.Context) bool {
	xr := ctx.Request().Header.Get("X-Requested-With")
	if xr != "" && strings.Index(xr, "XMLHttpRequest") != -1 {
		return true
	}
	accept := ctx.Request().Header.Get("Accept")
	if strings.Index(accept, "/json") != -1 || strings.Index(accept, "/+json") != -1 {
		return true
	}
	return false
}
