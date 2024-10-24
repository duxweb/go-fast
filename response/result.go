package response

import (
	"bytes"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/views"
	"github.com/go-errors/errors"
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, app string, name string, bind any, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}

	if bind == nil {
		bind = make(map[string]any)
	}

	if views.Views[app] == nil {
		return errors.New("tpl app not found")
	}

	logger.Log("tpl").Debug("name", name)

	buf := new(bytes.Buffer)
	err := views.Views[app].Render(buf, name, bind)
	if err != nil {
		return err
	}

	ctx.Set("content-type", fiber.MIMETextHTMLCharsetUTF8)
	return ctx.Status(statusCode).Send(buf.Bytes())
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
