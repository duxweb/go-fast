package response

import (
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, name string, bind any, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	return ctx.Render(statusCode, name, bind)
}

type Data struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"ok"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta"`
}

func Send(ctx echo.Context, data Data, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	if data.Meta == nil {
		data.Meta = echo.Map{}
	}
	data.Code = statusCode
	return ctx.JSON(statusCode, data)
}
