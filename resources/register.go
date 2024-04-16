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

func Register(files []*annotation.File) {

	for _, file := range files {
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

		appName := info.Params["app"].(string)
		resName := info.Params["name"].(string)

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
			if routeGroup != nil {
				routeGroup.Add(params["method"].(string), params["route"].(string), item.Func.(echo.HandlerFunc), params["name"].(string))
			}
			if permissionGroup != nil {
				permissionGroup.Add(params["name"].(string))
			}
		}

	}

}
