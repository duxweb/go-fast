package resources

import (
	"github.com/duxweb/go-fast/annotation"
	"github.com/duxweb/go-fast/permission"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
)

var Resources = map[string]*ResourceData{}

func Set(name string, data *ResourceData) {
	Resources[name] = data
}

func Get(name string) *ResourceData {
	return Resources[name]
}

func Register() {
	for _, file := range annotation.Annotations {
		// 获取资源数据
		var info *annotation.Annotation
		for _, item := range file.Annotations {
			if item.Name != "Resource" {
				continue
			}
			info = item
		}
		if info == nil {
			continue
		}

		appName, ok := info.Params["app"].(string)
		if !ok {
			panic("resource app not set: " + file.Name)
		}
		resName, ok := info.Params["name"].(string)
		if !ok {
			panic("resource name not set: " + file.Name)
		}

		// 设置路由组
		routeData := route.Get(appName)
		var routeGroup *route.RouterData
		if routeData != nil {
			routeGroup = routeData.Group(info.Params["route"].(string), resName)
		}
		permissionData := permission.Get(appName)
		var permissionGroup *permission.PermissionData
		// 设置权限组
		if permissionData != nil {
			permissionGroup = permissionData.Group(resName, 0)

		}
		// 设置资源动作
		for _, item := range file.Annotations {
			if item.Name != "Action" {
				continue
			}
			params := item.Params

			method, ok := params["method"].(string)
			if !ok {
				panic("action method not set: " + file.Name)
			}
			match, ok := params["route"].(string)
			if !ok {
				panic("action route not set: " + file.Name)
			}
			name, ok := params["name"].(string)
			if !ok {
				panic("action name not set: " + file.Name)
			}
			function, ok := item.Func.(func(ctx echo.Context) error)
			if !ok {
				panic("action func not set: " + file.Name)
			}

			if routeGroup != nil {
				routeGroup.Add(method, match, function, name)
			}
			if permissionGroup != nil {
				permissionGroup.Add(name)
			}
		}

	}

}
