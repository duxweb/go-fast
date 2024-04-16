package route

import (
	"github.com/duxweb/go-fast/annotation"
	"github.com/labstack/echo/v4"
)

var Routes = map[string]*RouterData{}

func Set(name string, route *RouterData) *RouterData {
	Routes[name] = route
	return route
}

func Get(name string) *RouterData {
	return Routes[name]
}

func Register(files []*annotation.File) {

	for _, file := range files {
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
			appName := info.Params["app"].(string)
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
			if routeGroup != nil {
				routeGroup.Add(params["method"].(string), params["route"].(string), item.Func.(echo.HandlerFunc), params["name"].(string))
			} else if routeData != nil {
				routeData.Add(params["method"].(string), params["route"].(string), item.Func.(echo.HandlerFunc), params["name"].(string))
			}
		}

	}

}
