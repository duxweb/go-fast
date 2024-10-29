package response

import (
	"github.com/duxweb/go-fast/i18n"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, app string, name string, bind any, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}

	ctx.Set("tpl", app)
	return ctx.Render(statusCode, name, bind)
}

type Data struct {
	Code        int    `json:"code" example:"200"`
	Message     string `json:"message" example:"ok"`
	MessageLang string `json:"-"`
	Data        any    `json:"data"`
	Meta        any    `json:"meta"`
}

func Send(ctx echo.Context, data Data, code ...int) error {
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
		data.Meta = echo.Map{}
	}
	if data.Code == 0 {
		data.Code = statusCode
	}
	return ctx.JSON(statusCode, data)
}
