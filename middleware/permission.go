package middleware

import (
	duxAuth "github.com/duxweb/go-fast/auth"
	duxPermission "github.com/duxweb/go-fast/permission"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type PermissionFun func(id string) ([]string, error)

func PermissionMiddleware(permission PermissionFun) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, ok := c.Get("auth").(*duxAuth.JwtClaims)
			if !ok {
				return response.BusinessError("Permissions must be authorized by the user after", 500)
			}
			routeName := route.GetRouteName(c)
			if routeName == "" {
				return next(c)
			}

			data := duxPermission.Get("admin").Get(c)
			permissions := lo.Map[map[string]any, string](data, func(item map[string]any, index int) string {
				return cast.ToString(item["name"])
			})
			c.Set("permissions", permissions)

			userPermission, err := permission(auth.ID)
			if err != nil {
				return err
			}
			c.Set("userPermissions", userPermission)

			err = Can(c, routeName)
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}

func Can(c echo.Context, name string) error {
	userPermission := c.Get("userPermissions").([]string)
	permission := c.Get("permissions").([]string)

	if len(userPermission) == 0 || len(permission) == 0 {
		return nil
	}

	if lo.IndexOf[string](permission, name) == -1 {
		return nil
	}

	if lo.IndexOf[string](userPermission, name) != -1 {
		return nil
	}

	return echo.ErrForbidden
}
