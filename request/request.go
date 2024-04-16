package request

import (
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
)

// RequestParser 请求解析验证
func RequestParser(ctx echo.Context, params any) error {
	var err error
	if err = ctx.Bind(params); err != nil {
		return err
	}
	err = validator.Validator().Struct(params)
	if err = validator.ProcessError(params, err); err != nil {
		return err
	}
	return nil
}
