package route

import (
	"github.com/labstack/echo/v4"
)

func AppMiddleware(route *RouterItem) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("route", route)
			return next(c)
		}
	}
}
