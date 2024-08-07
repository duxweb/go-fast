package middleware

import (
	"github.com/duxweb/go-fast/helper"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func VisitorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method

		if method != "GET" {
			return next(c)
		}

		path := c.Path()
		hasInstallLock := helper.IsExist("./data/install.lock")
		pathContainsInstall := strings.Contains(path, "/install")

		if !hasInstallLock && !pathContainsInstall {
			err := c.Redirect(http.StatusFound, "/install")
			return err
		}

		if hasInstallLock && pathContainsInstall {
			err := c.Redirect(http.StatusFound, "/")
			return err
		}

		err := helper.VisitIncrement(c, "common", 0, "web", "")
		if err != nil {
			return err
		}

		return next(c)
	}
}
