package middleware

import (
	"github.com/duxweb/go-fast/helper"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

func VisitorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()

		if method != "GET" {
			return c.Next()
		}

		path := c.Path()
		hasInstallLock := helper.IsExist("./data/install.lock")
		pathContainsInstall := strings.Contains(path, "/install")

		if !hasInstallLock && !pathContainsInstall {
			err := c.Redirect("/install", http.StatusFound)
			return err
		}

		if hasInstallLock && pathContainsInstall {
			err := c.Redirect("/", http.StatusFound)
			return err
		}

		err := helper.VisitIncrement(c, "common", 0, "web", "")
		if err != nil {
			return err
		}

		return c.Next()
	}
}
