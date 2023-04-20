package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/dashboard"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/handlers"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/color"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
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
			var msg string
			if e, ok := err.(*fiber.Error); ok {
				// http error
				code = e.Code
				msg = e.Message
			} else {
				// Other error
				marshal, _ := json.Marshal(eris.ToJSON(eris.Wrapf(err, "error"), true))
				logger.Log().Error().RawJSON("stack", marshal).Msg(err.Error())
				msg = lo.Ternary[string](global.DebugMsg == "", "business is busy, please try again", global.DebugMsg)
			}

			// Asynchronous request
			if ctx.Is("json") || ctx.XHR() {
				ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return ctx.Status(code).JSON(handlers.New(code, err.Error()))
			}

			// Web request
			if code == http.StatusNotFound {
				ctx.Status(code).Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				return ctx.Render("template/404", fiber.Map{})
			} else {
				ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				return ctx.Render("template/error", fiber.Map{
					"code":    code,
					"message": msg,
				})
			}
		},
		Views: views.Views,
	})

	// Exception recovery processing
	global.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			marshal, _ := json.Marshal(eris.ToJSON(eris.New("panic"), true))
			logger.Log().Error().Interface("message", e).RawJSON("stack", marshal)
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

	// Dashboard
	dashboard.Init()

	// Default route
	global.App.Get("/", func(c *fiber.Ctx) error {
		return c.Render("template/welcome", fiber.Map{})
	})

	// Websocket route
	global.App.Get("/ws", func(c *fiber.Ctx) error {
		data, err := auth.NewJWT().ParsingToken(c.Get("Authorization"))
		if err != nil {
			return err
		}
		return websocket.Socket.Handler(gocast.Str(data["sub"]), gocast.Str(data["id"]))(c)
	})

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
