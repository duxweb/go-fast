package route

import (
	"github.com/duxweb/go-fast/global"
	"github.com/gofiber/fiber/v2"
)

type RouterData struct {
	title      string
	prefix     string
	permission bool
	data       []*RouterItem
	group      []*RouterData
	router     fiber.Router
}

type RouterItem struct {
	title  string
	method string
	path   string
	name   string
}

func New(prefix string, middle ...fiber.Handler) *RouterData {
	return &RouterData{
		router: global.App.Group(prefix, middle...),
	}
}

func (t *RouterData) Group(prefix string, title string, middle ...fiber.Handler) *RouterData {
	group := &RouterData{
		title:  title,
		prefix: prefix,
		router: t.router.Group(prefix, middle...),
	}
	t.group = append(t.group, group)
	return group
}

func (t *RouterData) Permission() *RouterData {
	t.permission = true
	return t
}

func (t *RouterData) Router() fiber.Router {
	return t.router
}

func (t *RouterData) Get(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("GET", path, handler, title, name)
}

func (t *RouterData) Head(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("HEAD", path, handler, title, name)
}

func (t *RouterData) Post(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("POST", path, handler, title, name)
}

func (t *RouterData) Put(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("PUT", path, handler, title, name)
}

func (t *RouterData) Delete(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("DELETE", path, handler, title, name)
}

func (t *RouterData) Connect(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("CONNECT", path, handler, title, name)
}

func (t *RouterData) Options(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("OPTIONS", path, handler, title, name)
}

func (t *RouterData) Trace(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("TRACE", path, handler, title, name)
}

func (t *RouterData) Patch(path string, handler fiber.Handler, title string, name string) fiber.Router {
	return t.Add("PATH", path, handler, title, name)
}

func (t *RouterData) Add(method string, path string, handler fiber.Handler, title string, name string) fiber.Router {
	item := RouterItem{
		title:  title,
		method: method,
		path:   path,
		name:   name,
	}
	t.data = append(t.data, &item)
	return t.router.Add(method, path, handler).Name(item.name)
}

func (t *RouterData) ParseTree(prefix string) any {
	var all []any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"title":  datum.title,
			"name":   datum.name,
			"method": datum.method,
			"path":   prefix + datum.path,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		all = append(all, item.ParseTree(gpath))
	}
	if t.title == "" {
		return all
	}
	return map[string]any{
		"title": t.title,
		"path":  prefix,
		"data":  all,
	}
}

func (t *RouterData) ParseData(prefix string) []map[string]any {
	var all []map[string]any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"title":  datum.title,
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
