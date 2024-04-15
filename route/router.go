package route

import (
	"github.com/duxweb/go-fast/global"
	"github.com/labstack/echo/v4"
)

type RouterData struct {
	name   string
	prefix string
	data   []*RouterItem
	group  []*RouterData
	router *echo.Group
}

type RouterItem struct {
	method string
	path   string
	name   string
}

func New(prefix string, middle ...echo.MiddlewareFunc) *RouterData {
	return &RouterData{
		router: global.App.Group(prefix, middle...),
	}
}

func (t *RouterData) Group(prefix string, name string, middle ...echo.MiddlewareFunc) *RouterData {
	group := &RouterData{
		prefix: prefix,
		router: t.router.Group(prefix, middle...),
		name:   name,
	}
	t.group = append(t.group, group)
	return group
}

func (t *RouterData) Router() *echo.Group {
	return t.router
}

func (t *RouterData) Get(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("GET", path, handler, name)
}

func (t *RouterData) Head(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("HEAD", path, handler, name)
}

func (t *RouterData) Post(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("POST", path, handler, name)
}

func (t *RouterData) Put(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("PUT", path, handler, name)
}

func (t *RouterData) Delete(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("DELETE", path, handler, name)
}

func (t *RouterData) Connect(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("CONNECT", path, handler, name)
}

func (t *RouterData) Options(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("OPTIONS", path, handler, name)
}

func (t *RouterData) Trace(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("TRACE", path, handler, name)
}

func (t *RouterData) Patch(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("PATH", path, handler, name)
}

func (t *RouterData) Any(path string, handler echo.HandlerFunc, name string) *echo.Route {
	return t.Add("ANY", path, handler, name)
}

func (t *RouterData) Add(method string, path string, handler echo.HandlerFunc, name string) *echo.Route {
	item := RouterItem{
		method: method,
		path:   path,
		name:   name,
	}
	t.data = append(t.data, &item)
	r := t.router.Add(method, path, handler)
	r.Name = item.name
	return r
}

func (t *RouterData) ParseTree(prefix string) any {
	var all []any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"name":   datum.name,
			"method": datum.method,
			"path":   prefix + datum.path,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		all = append(all, item.ParseTree(gpath))
	}
	return map[string]any{
		"path": prefix,
		"data": all,
	}
}

func (t *RouterData) ParseData(prefix string) []map[string]any {
	var all []map[string]any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"name":   datum.name,
			"method": datum.method,
			"path":   prefix + datum.path,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		data := item.ParseData(gpath)
		all = append(all, data...)
	}
	return all
}
