package route

import (
	"github.com/duxweb/go-fast/annotation"
	"github.com/gofiber/fiber/v2"
)

var Routes = map[string]*RouterData{}

func Set(name string, route *RouterData) *RouterData {
	Routes[name] = route
	return route
}

func Get(name string) *RouterData {
	return Routes[name]
}

func Register() {
	for _, file := range annotation.Annotations {
		// 获取资源数据
		var info *annotation.Annotation
		for _, item := range file.Annotations {
			if item.Name != "RouteGroup" {
				continue
			}
			info = item
		}

		// 设置路由组
		var routeGroup *RouterData
		var routeData *RouterData
		if info != nil {
			appName, ok := info.Params["app"].(string)
			if !ok {
				panic("the routing group does not have app parameters set: " + file.Name)
			}
			resName := info.Params["name"].(string)
			routeData = Get(appName)
			if routeData != nil {
				routeGroup = routeData.Group(info.Params["route"].(string), resName)
			}
		}

		// 设置路由
		for _, item := range file.Annotations {
			if item.Name != "Route" {
				continue
			}
			params := item.Params

			method, ok := params["method"].(string)
			if !ok {
				panic("routing method not set: " + file.Name)
			}
			match, ok := params["route"].(string)
			if !ok {
				panic("routing route not set: " + file.Name)
			}
			name, ok := params["name"].(string)
			if !ok {
				panic("routing name not set: " + file.Name)
			}
			function, ok := item.Func.(func(ctx *fiber.Ctx) error)
			if !ok {
				panic("routing func not set: " + file.Name)
			}
			if routeGroup != nil {
				routeGroup.Add(method, match, function, name)
			} else if routeData != nil {
				routeData.Add(method, match, function, name)
			}
		}

	}

}
