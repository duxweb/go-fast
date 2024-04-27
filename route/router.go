package route

import (
	"github.com/duxweb/go-fast/global"
	"github.com/labstack/echo/v4"
)

type RouterData struct {
	Name        string
	Prefix      string
	Data        []*RouterItem
	Groups      []*RouterData
	GroupRouter *echo.Group
}

type RouterItem struct {
	Method string
	Path   string
	Name   string
}

func New(prefix string, middle ...echo.MiddlewareFunc) *RouterData {
	return &RouterData{
		Prefix:      prefix,
		GroupRouter: global.App.Group(prefix, middle...),
	}
}

func (t *RouterData) Group(prefix string, name string, middle ...echo.MiddlewareFunc) *RouterData {
	group := &RouterData{
		Prefix:      prefix,
		GroupRouter: t.GroupRouter.Group(prefix, middle...),
		Name:        name,
	}
	t.Groups = append(t.Groups, group)
	return group
}

func (t *RouterData) Router() *echo.Group {
	return t.GroupRouter
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
		Method: method,
		Path:   path,
		Name:   name,
	}
	t.Data = append(t.Data, &item)
	r := t.GroupRouter.Add(method, path, handler, AppMiddleware(&item))
	r.Name = item.Name
	return r
}

func (t *RouterData) ParseTree(prefix string) any {
	var all []any
	for _, datum := range t.Data {
		all = append(all, map[string]any{
			"name":   datum.Name,
			"method": datum.Method,
			"path":   prefix + datum.Path,
		})
	}
	for _, item := range t.Groups {
		gpath := prefix + item.Prefix
		all = append(all, item.ParseTree(gpath))
	}
	return map[string]any{
		"path": prefix,
		"data": all,
	}
}

func (t *RouterData) ParseData(prefix string) []map[string]any {
	var all []map[string]any
	for _, datum := range t.Data {
		all = append(all, map[string]any{
			"name":   datum.Name,
			"method": datum.Method,
			"path":   prefix + datum.Path,
		})
	}
	for _, item := range t.Groups {
		gpath := prefix + item.Prefix
		data := item.ParseData(gpath)
		all = append(all, data...)
	}
	return all
}
