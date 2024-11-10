package web

import "github.com/labstack/echo/v4"

func CacheHandler(value string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, value)
			return next(c)
		}
	}
}
