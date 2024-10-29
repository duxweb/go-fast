package middleware

import (
	duxAuth "github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
)

type PermissionFun func(id string) (map[string]bool, error)

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

			permissions, err := permission(auth.ID)
			if err != nil {
				return err
			}

			err = Can(permissions, routeName)
			if err != nil {
				return err
			}
			return next(c)
		}
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
		return echo.ErrForbidden
	}
	return nil
}
