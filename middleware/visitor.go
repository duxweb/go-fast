package middleware

import (
	"github.com/duxweb/go-fast/helper"
	"github.com/labstack/echo/v4"
)

func VisitorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method

		if method != "GET" {
			return next(c)
		}

		err := helper.VisitIncrement(c, "common", 0, "web", "")
		if err != nil {
			return err
		}

		return next(c)
	}
}
