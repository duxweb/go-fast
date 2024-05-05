package resources

import (
	"github.com/duxweb/go-fast/menu"
	"github.com/duxweb/go-fast/middleware"
	"github.com/duxweb/go-fast/permission"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
)

type ResourceData struct {
	name           string
	path           string
	authMiddleware []echo.MiddlewareFunc
	middleware     []echo.MiddlewareFunc
	permission     middleware.PermissionFun
	operate        bool
}

func New(name string, path string) *ResourceData {
	return &ResourceData{
		name: name,
		path: path,
	}

}

func (t *ResourceData) AddMiddleware(middle ...echo.MiddlewareFunc) *ResourceData {
	t.authMiddleware = append(t.middleware, middle...)
	return t
}

func (t *ResourceData) AddAuthMiddleware(middle ...echo.MiddlewareFunc) *ResourceData {
	t.authMiddleware = append(t.authMiddleware, middle...)
	return t
}

func (t *ResourceData) GetMiddleware() []echo.MiddlewareFunc {
	return t.middleware
}

func (t *ResourceData) GetAuthMiddleware() []echo.MiddlewareFunc {
	return t.authMiddleware
}

func (t *ResourceData) GetAllMiddleware() []echo.MiddlewareFunc {
	return append(t.middleware, t.authMiddleware...)
}

func (t *ResourceData) SetPermission(getPermission middleware.PermissionFun) *ResourceData {
	t.permission = getPermission
	return t
}

func (t *ResourceData) SetOperate(status bool) *ResourceData {
	t.operate = status
	return t
}

func (t *ResourceData) run() *ResourceData {

	middle := []echo.MiddlewareFunc{
		middleware.AuthMiddleware("admin"),
	}
	if t.permission != nil {
		middle = append(middle, middleware.PermissionMiddleware(t.permission))
	}
	if t.operate {
		middle = append(middle, middleware.OperateMiddleware(t.name))
	}
	middle = append(middle, t.GetAllMiddleware()...)

	route.Set(t.name, route.New(t.path, middle...))
	permission.Set(t.name, permission.New())
	menu.Set(t.name, menu.New(t.path))
	return t
}
