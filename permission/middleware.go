package permission

import (
	duxAuth "github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"sync"
)

func PermissionMiddleware(app string, getPermission func(id int64) []string) echo.MiddlewareFunc {
	var doOnce sync.Once
	var permissions []string
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			auth, ok := c.Get("auth").(duxAuth.JwtClaims)
			if !ok {
				return echo.ErrUnauthorized
			}
			doOnce.Do(func() {
				permissions = Get(app).GetData()
			})
			routeItem := c.Get("route").(*route.RouterItem)
			routeName := routeItem.Name
			if routeName == "" || lo.IndexOf[string](permissions, routeName) == -1 {
				return next(c)
			}
			permissions = getPermission(cast.ToInt64(auth.ID))
			if len(permissions) > 0 && lo.IndexOf[string](permissions, routeName) == -1 {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}
