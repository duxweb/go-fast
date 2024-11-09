package web

import (
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/response"
	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		result := response.Data{
			Code: http.StatusInternalServerError,
		}
		var e *echo.HTTPError
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
				logger.Log().Error("core", "err", err, slog.String("stack", string(debug.Stack())))
				result.Message = err.Error()
			}

			result.Message = lo.Ternary[string](!global.Debug, i18n.Get(c, "common.error.errorMessage"), result.Message)
		}

		if isAsync(c) {
			err = response.Send(c, result, result.Code)
			if err != nil {
				logger.Log().Error("err", err)
			}
			return
		}

		c.Set("tpl", "app")

		if result.Code == http.StatusNotFound {
			err = c.Render(http.StatusNotFound, "template/404.html", nil)
		} else {
			err = c.Render(http.StatusInternalServerError, "template/error.html", map[string]any{
				"code":    result.Code,
				"message": result.Message,
			})
		}
		if err != nil {
			logger.Log().Error("err", err)
		}
	}
}

func isAsync(ctx echo.Context) bool {
	xr := ctx.Request().Header.Get("X-Requested-With")
	if xr != "" && strings.Index(xr, "XMLHttpRequest") != -1 {
		return true
	}
	accept := ctx.Request().Header.Get("Accept")
	if strings.Index(accept, "/json") != -1 || strings.Index(accept, "/+json") != -1 {
		return true
	}
	return false
}
