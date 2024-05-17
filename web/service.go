package web

import (
	"context"
	"fmt"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/websocket"
	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cast"
	"io/fs"
	"log/slog"
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
	global.App.Debug = global.Debug
	global.App.Renderer = views.Render()
	global.App.HideBanner = true
	global.App.HidePort = true

	global.App.Logger = EchoLoggerHeadAdaptor()

	// 注册异常处理
	global.App.HTTPErrorHandler = func(err error, c echo.Context) {
		result := response.Data{
			Code: http.StatusInternalServerError,
		}
		var e *echo.HTTPError
		if errors.As(err, &e) {
			// http error
			result.Code = e.Code
			result.Message = cast.ToString(e.Message)
		} else {
			var exceptions *errors.Error
			var validator *response.ValidatorData
			if errors.As(err, &exceptions) {
				stacks := exceptions.StackFrames()
				logger.Log().Error("core", "err", err,
					slog.String("file", lo.Ternary[string](len(stacks) > 0, stacks[0].File+":"+cast.ToString(stacks[0].LineNumber), "")),
					slog.Any("stack", lo.Map[errors.StackFrame, map[string]any](stacks, func(item errors.StackFrame, index int) map[string]any {
						return map[string]any{
							"file": item.File + ":" + cast.ToString(item.LineNumber),
							"func": item.Name,
						}
					})),
				)
			} else if errors.As(err, &validator) {
				result.Code = validator.Code
				result.Data = validator.Data
				result.Message = validator.Message
			} else {
				logger.Log().Error("core", "err", err)
				result.Message = err.Error()
			}

			result.Message = lo.Ternary[string](!global.Debug, i18n.Trans.Get("common.error.errorMessage"), result.Message)
		}

		if isAsync(c) {
			err = response.Send(c, result, result.Code)
			if err != nil {
				logger.Log().Error("err", err)
			}
			return
		}

		c.Set("tpl", "app")

		if result.Code == http.StatusNotFound {
			err = c.Render(http.StatusNotFound, "404.gohtml", nil)
		} else {
			err = c.Render(http.StatusInternalServerError, "error.gohtml", map[string]any{
				"code":    result.Code,
				"message": result.Message,
			})
		}
		if err != nil {
			logger.Log().Error("err", err)
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

	// 注册公共目录
	if global.StaticFs != nil {
		entries, _ := fs.ReadDir(echo.MustSubFS(*global.StaticFs, "static"), ".")
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				global.App.StaticFS("/"+name, echo.MustSubFS(*global.StaticFs, "static/"+name))
			}
		}
	}

	// 注册静态路由
	global.App.Static("/", "./public")

	// 请求日志
	global.App.Use(middleware.RequestID())
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

			var level slog.Level
			attr := []slog.Attr{
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("ip", v.RemoteIP),
				slog.Duration("latency", v.Latency),
				slog.String("id", v.RequestID),
			}

			if v.Error != nil {
				level = slog.LevelError
				attr = append(attr, slog.Attr{Key: "err", Value: slog.StringValue(v.Error.Error())})
			} else {
				level = lo.Ternary[slog.Level](v.Latency > 1*time.Second, slog.LevelWarn, slog.LevelInfo)
			}

			logger.Log("request").LogAttrs(
				context.Background(),
				level,
				"request",
				attr...,
			)

			return nil
		},
	}))

	timeout := 60 * time.Second
	if config.IsLoad("use") {
		timeout = config.Load("use").GetDuration("server.timeout") * time.Second
	}

	global.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: timeout,
	}))
}

func Start() {

	global.App.GET("/", func(c echo.Context) error {
		c.Set("tpl", "app")
		return c.Render(http.StatusOK, "welcome.gohtml", nil)
	})

	// websocket
	global.App.GET("/ws", func(c echo.Context) error {
		token := c.QueryParam("token")
		app := c.QueryParam("app")
		if token == "" {
			logger.Log("websocket").Debug("Token Not Found", slog.String("token", token))
			return response.Send(c, response.Data{
				Message: "token does not exist",
			})
		}
		if app == "" {
			logger.Log("websocket").Debug("App Not Found", slog.String("token", token))
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

	port := "8900"
	if config.IsLoad("use") {
		port = config.Load("use").GetString("server.port")
	}

	banner()
	global.BootTime = time.Now()
	color.Println("⇨ <green>Server start http://localhost:" + port + "</>")

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
