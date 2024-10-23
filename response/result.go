package response

import (
	"github.com/duxweb/go-fast/i18n"
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, app string, name string, bind any, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	return ctx.Status(statusCode).Render(name, bind, app)
}

type Data struct {
	Code        int    `json:"code" example:"200"`
	Message     string `json:"message" example:"ok"`
	MessageLang string `json:"-"`
	Data        any    `json:"data"`
	Meta        any    `json:"meta"`
}

func Send(ctx *fiber.Ctx, data Data, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	if data.Message == "" {
		data.Message = "ok"
	}
	if data.MessageLang != "" {
		data.Message = i18n.Trans.Get(data.MessageLang)
	}
	if data.Meta == nil {
		data.Meta = fiber.Map{}
	}
	if data.Code == 0 {
		data.Code = statusCode
	}
	return ctx.Status(statusCode).JSON(data)
}
