package middleware

import (
	duxAuth "github.com/duxweb/go-fast/auth"
	"github.com/gofiber/fiber/v2"
)

type PermissionFun func(id string) (map[string]bool, error)

func PermissionMiddleware(permission PermissionFun) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, ok := c.Locals("auth").(*duxAuth.JwtClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}
		routeName := c.Route().Name
		if routeName == "" {
			return c.Next()
		}

		permissions, err := permission(auth.ID)
		if err != nil {
			return err
		}

		err = Can(permissions, routeName)
		if err != nil {
			return err
		}
		return c.Next()
	}
}

func Can(permissions map[string]bool, name string) error {
	if len(permissions) == 0 {
		return nil
	}
	is, ok := permissions[name]
	if !ok {
		return nil
	}
	if !is {
		return fiber.ErrForbidden
	}
	return nil
}
