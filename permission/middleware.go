package permission

import (
	"github.com/demdxx/gocast/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"sync"
)

// Middleware permission
func Middleware(app string, callback func(id int64) []string) fiber.Handler {
	var doOnce sync.Once
	var permissions []string
	return func(c *fiber.Ctx) error {
		auth, ok := c.Locals("auth").(map[string]any)
		if !ok {
			return fiber.ErrUnauthorized
		}
		doOnce.Do(func() {
			permissions = Get(app).GetData()
		})
		routeName := c.Route().Name
		if routeName == "" || lo.IndexOf[string](permissions, routeName) == -1 {
			return c.Next()
		}
		permissions = callback(gocast.Int64(auth["id"]))
		if len(permissions) > 0 && lo.IndexOf[string](permissions, routeName) == -1 {
			return fiber.ErrForbidden
		}
		return c.Next()
	}
}
