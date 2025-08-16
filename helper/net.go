package helper

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
)

func NetIsAjax(ctx echo.Context) bool {
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

func SetContextValue[K any, V any](ctx context.Context, key K, value V) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContextValue[K any, V any](ctx context.Context, key K) V {
	value := ctx.Value(key)
	if val, ok := value.(V); ok {
		return val
	}
	var zero V
	return zero
}

func EchoSetContextValue[K any, V any](c echo.Context, key K, value V) {
	ctx := context.WithValue(c.Request().Context(), key, value)
	c.SetRequest(c.Request().WithContext(ctx))
}

func EchoGetContextValue[K any, V any](c echo.Context, key K) V {

	value := c.Request().Context().Value(key)

	if val, ok := value.(V); ok {
		return val
	}

	var zero V
	return zero
}
