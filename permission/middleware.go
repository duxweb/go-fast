package permission

import (
	duxAuth "github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
)

func PermissionMiddleware(getPermission func(id string) map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, ok := c.Get("auth").(duxAuth.JwtClaims)
			if !ok {
				return echo.ErrUnauthorized
			}
			routeItem := c.Get("route").(*route.RouterItem)
			routeName := routeItem.Name
			if routeName == "" {
				return next(c)
			}
			permissions := getPermission(auth.ID)

			err := Can(permissions, routeName)
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}
