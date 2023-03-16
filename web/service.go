package web

import (
	"errors"
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/handlers"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/color"
	"github.com/samber/lo"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

func Init() {

	proxyHeader := config.Get("app").GetString("app.proxyHeader")
	global.App = fiber.New(fiber.Config{
		AppName:               "DuxGO",
		Prefork:               false,
		CaseSensitive:         false,
		StrictRouting:         false,
		EnablePrintRoutes:     false,
		DisableStartupMessage: true,
		ProxyHeader:           lo.Ternary[string](proxyHeader != "", proxyHeader, "X-Real-IP"),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var msg any
			if e, ok := err.(*handlers.CoreError); ok {
				// program error
				msg = e.Message
			} else if e, ok := err.(*fiber.Error); ok {
				// http error
				code = e.Code
				msg = e.Message
			} else {
				// Other error
				msg = err.Error()
				logger.Log().Error().Bytes("body", ctx.Body()).Err(err).Msg("error")
			}

			// Asynchronous request
			if ctx.Is("json") || ctx.XHR() {
				ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return ctx.Status(code).JSON(handlers.New(code, err.Error()))
			}

			// Web request
			if code == http.StatusNotFound {
				ctx.Status(code).Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				return views.FrameTpl.ExecuteTemplate(ctx.Response().BodyWriter(), "404.gohtml", nil)
			} else {
				ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				err = views.FrameTpl.ExecuteTemplate(ctx.Response().BodyWriter(), "error.gohtml", fiber.Map{
					"code":    code,
					"message": msg,
				})
				if err != nil {
					logger.Log().Error().Err(err).Send()
				}
				return nil
			}
		},
		Views: views.Views,
	})

	// Exception recovery processing
	global.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.Log().Error().Interface("err", e).Bytes("stack", debug.Stack()).Send()
		},
	}))

	// Cors process across domains
	global.App.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "*",
		ExposeHeaders: "*",
	}))

	// Request log
	global.App.Use(fiberLogger())

	// Registering websocket
	websocket.Init()
}

func Start() {
	port := config.Get("app").GetString("server.port")
	banner()
	global.BootTime = time.Now()
	color.Println("⇨ <green>Server start http://0.0.0.0:" + port + "</>")
	go func() {
		err := global.App.Listen(":" + port)
		if errors.Is(err, http.ErrServerClosed) {
			color.Println("⇨ <red>Server closed</>")
			return
		}
		if err != nil {
			logger.Log().Error().Err(err).Msg("web")
		}
	}()
}

func banner() {
	debugBool := config.Get("app").GetBool("server.debug")

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
		Value: lo.Ternary[string](debugBool, "enabled", "disabled"),
	})
	sysMaps = append(sysMaps, item{
		Name:  "PID",
		Value: os.Getpid(),
	})
	sysMaps = append(sysMaps, item{
		Name:  "Routes",
		Value: len(global.App.Stack()),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	color.Println(banner)
}
