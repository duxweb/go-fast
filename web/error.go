package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/duxweb/go-fast/v2/errors"
	"github.com/duxweb/go-fast/v2/helper"
	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
	"github.com/spf13/cast"
)

type TError struct {
	Code    int
	Message string
	Data    any
}

// ErrorHandler 错误处理
// ErrorHandler error handler
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {

		result := TError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}

		switch e := err.(type) {
		case *echo.HTTPError:
			result.Code = e.Code
			result.Message = cast.ToString(e.Message)
		case errors.HTTPError:
			result.Code = e.Code
			result.Message = e.Error()
			result.Data = e.Data
		case oops.OopsError:
			result.Message = e.Error()
			c.Logger().Error(e.Error(), slog.Any("error", err))
		default:
			result.Message = err.Error()
			c.Logger().Error(err.Error(), slog.Any("error", err))
		}

		if helper.NetIsAjax(c) {
			err = c.JSON(result.Code, result)
			if err != nil {
				c.Logger().Error(err.Error(), slog.Any("error", err))
			}
			return
		}

		err = c.Render(result.Code, "template/error.html", map[string]any{
			"title":   fmt.Sprintf("%d | %s", result.Code, result.Message),
			"code":    result.Code,
			"message": result.Message,
		})

		if err != nil {
			c.Logger().Error(err.Error(), slog.Any("error", err))
		}

	}
}
