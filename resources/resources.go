package resources

import (
	"github.com/duxweb/go-fast/permission"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
)

type ResourceData struct {
	name           string
	path           string
	authMiddleware []echo.MiddlewareFunc
	middleware     []echo.MiddlewareFunc
}

func New() *ResourceData {
	return &ResourceData{}
}

func (t *ResourceData) addMiddleware(middle ...echo.MiddlewareFunc) *ResourceData {
	t.authMiddleware = append(t.middleware, middle...)
	return t
}

func (t *ResourceData) addAuthMiddleware(middle ...echo.MiddlewareFunc) *ResourceData {
	t.authMiddleware = append(t.authMiddleware, middle...)
	return t
}

func (t *ResourceData) getMiddleware() []echo.MiddlewareFunc {
	return t.middleware
}

func (t *ResourceData) getAuthMiddleware() []echo.MiddlewareFunc {
	return t.authMiddleware
}

func (t *ResourceData) getAllMiddleware() []echo.MiddlewareFunc {
	return append(t.middleware, t.authMiddleware...)
}

func (t *ResourceData) run() {
	route.Set(t.name, route.New(t.path))
	permission.Set(t.name, permission.New())
}
