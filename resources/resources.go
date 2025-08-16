package resources

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/middleware"
	"github.com/duxweb/go-fast/v2/permission"
	"github.com/duxweb/go-fast/v2/route"
)

type ResourceData struct {
	name           string
	path           string
	authMiddleware huma.Middlewares
	middleware     huma.Middlewares
	permission     middleware.PermissionFun
	operate        bool
}

func New(name string, path string) *ResourceData {
	return &ResourceData{
		name: name,
		path: path,
	}
}

func (t *ResourceData) AddMiddleware(middle huma.Middlewares) *ResourceData {
	t.middleware = append(t.middleware, middle...)
	return t
}

func (t *ResourceData) AddAuthMiddleware(middle huma.Middlewares) *ResourceData {
	t.authMiddleware = append(t.authMiddleware, middle...)
	return t
}

func (t *ResourceData) GetMiddleware() huma.Middlewares {
	return t.middleware
}

func (t *ResourceData) GetAuthMiddleware() huma.Middlewares {
	return t.authMiddleware
}

func (t *ResourceData) GetAllMiddleware() huma.Middlewares {
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

	middle := []func(huma.Context, func(huma.Context)){
		middleware.AuthMiddleware(t.name),
	}
	if t.permission != nil {
		middle = append(middle, middleware.PermissionMiddleware(t.permission))
	}
	if t.operate {
		middle = append(middle, middleware.OperateMiddleware(t.name))
	}
	middle = append(middle, t.GetAllMiddleware()...)

	route.SetRouter(t.name, route.New(t.name, t.path, middle...))
	permission.Set(t.name, permission.New())
	return t
}
