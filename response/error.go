package response

import (
	"context"

	"github.com/duxweb/go-fast/i18n"
	"github.com/labstack/echo/v4"
)

func BusinessError(message any, code ...int) error {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}
	return echo.NewHTTPError(statusCode, message)
}

func GetEchoCtx(ctx context.Context) echo.Context {
	return ctx.Value("echo").(echo.Context)
}

func BusinessLangError(ctx echo.Context, message string, code ...int) error {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}
	return echo.NewHTTPError(statusCode, i18n.Get(ctx, message))
}

type ValidatorData Data

func (t *ValidatorData) Error() string {
	return t.Message
}

func ValidatorError(message string, data any, code ...int) error {
	statusCode := 422
	if len(code) > 0 {
		statusCode = code[0]
	}
	return &ValidatorData{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
}
