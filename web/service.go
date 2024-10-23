package web

import (
	"fmt"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/websocket"
	"github.com/go-errors/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gookit/goutil/fsutil"
	"github.com/spf13/cast"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/logger"
	fiberlog "github.com/gofiber/fiber/v2/log"
	middleLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/color"
	"github.com/samber/lo"
)

func Init() {
	global.App = fiber.New(fiber.Config{
		AppName:   "Dux",
		Immutable: true,
		//EnablePrintRoutes:     global.Debug,
		DisableStartupMessage: true,
		Views:                 views.Fiber(),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			result := response.Data{
				Code: fiber.StatusInternalServerError,
			}
			var e *fiber.Error

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

			if isAsync(ctx) {
				return response.Send(ctx, result, result.Code)
			}
			ctx.Locals("tpl", "app")

			if result.Code == fiber.StatusNotFound {
				return ctx.Status(fiber.StatusNotFound).Render("404.gohtml", nil)
			} else {
				return ctx.Status(fiber.StatusInternalServerError).Render("error.gohtml", fiber.Map{
					"code":    result.Code,
					"message": result.Message,
				})
			}
		},
	})

	// 适配日志
	fiberlog.SetLogger(LoggerAdaptor())

	// 注册公共目录
	global.App.Static("/public", "./public")

	// 请求日志
	global.App.Use(requestid.New())
	global.App.Use(middleLogger.New(middleLogger.Config{
		Next: func(c *fiber.Ctx) bool {
			if strings.Contains(c.Path(), "/public/") {
				return true
			}
			return false
		},
	}))

	// 异常恢复
	global.App.Use(recover.New(recover.Config{
		Next:             nil,
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.Log().Error("panic", fmt.Sprintf("%v", e),
				slog.Any("stack", string(debug.Stack())),
			)
		},
	}))

	// 跨域处理
	global.App.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "*",
		ExposeHeaders: "*",
	}))

	// ETAG 缓存
	global.App.Use(etag.New())

	// 图表
	if fsutil.IsFile("./public/favicon.ico") {
		global.App.Use(favicon.New(favicon.Config{
			File: "./public/favicon.ico",
			URL:  "/favicon.ico",
		}))
	}

	// 注册静态路由
	if global.StaticFs != nil {
		global.App.Use("/static", filesystem.New(filesystem.Config{
			Root:       http.FS(global.StaticFs),
			PathPrefix: "static",
		}))
	}

	//timeout := 60 * time.Second
	//if config.IsLoad("use") {
	//	timeout = config.Load("use").GetDuration("server.timeout") * time.Second
	//}

	global.App.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).Render("welcome.gohtml", nil, "app")
	})

	// websocket
	global.App.Get("/ws", func(c *fiber.Ctx) error {
		token := c.Query("token")
		app := c.Query("app")
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
		handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("token", cast.ToString(token))
			err := websocket.Service.Websocket.HandleRequest(w, r)
			if err != nil {
				logger.Log("websocket").Error(err.Error())
			}
		})
		return handler(c)
	})
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
		err := global.App.Listen(":" + port)
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
		Name:  "Fiber",
		Value: fiber.Version,
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
		Value: len(global.App.GetRoutes(true)),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	color.Println(banner)
}

func isAsync(ctx *fiber.Ctx) bool {
	xr := ctx.GetRespHeader("X-Requested-With")
	if xr != "" && strings.Index(xr, "XMLHttpRequest") != -1 {
		return true
	}
	accept := ctx.GetRespHeader("Accept")
	if strings.Index(accept, "/json") != -1 || strings.Index(accept, "/+json") != -1 {
		return true
	}
	return false
}
