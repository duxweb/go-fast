package resources

import (
	"github.com/duxweb/go-fast/action"
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

		resFunc, ok := info.Func.(func() action.Result)
		if !ok {
			panic("resource fun not set: " + file.Name)
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

		// 设置内置资源
		resFuncMap := resFunc()
		if resFuncMap["list"] != nil {
			routeGroup.Add("GET", "", resFuncMap["list"], resName+".list")
		}
		if resFuncMap["show"] != nil {
			routeGroup.Add("GET", "/{id:[0-9]+}", resFuncMap["show"], resName+".show")
		}
		if resFuncMap["create"] != nil {
			routeGroup.Add("POST", "", resFuncMap["create"], resName+".create")
		}
		if resFuncMap["edit"] != nil {
			routeGroup.Add("PUT", "/{id:[0-9]+}", resFuncMap["edit"], resName+".edit")
		}
		if resFuncMap["store"] != nil {
			routeGroup.Add("PATH", "/{id:[0-9]+}", resFuncMap["store"], resName+".store")
		}
		if resFuncMap["delete"] != nil {
			routeGroup.Add("DELETE", "/{id:[0-9]+}", resFuncMap["delete"], resName+".delete")
			routeGroup.Add("DELETE", "", resFuncMap["deleteMany"], resName+".deleteMany")
		}
		if resFuncMap["trash"] != nil {
			routeGroup.Add("PATH", "/{id:[0-9]+}/trash", resFuncMap["trash"], resName+".trash")
		}
		if resFuncMap["restore"] != nil {
			routeGroup.Add("PATH", "/{id:[0-9]+}/restore", resFuncMap["restore"], resName+".restore")
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
				routeGroup.Add(method, match, function, resName+"."+name)
			}
			if permissionGroup != nil {
				permissionGroup.Add(name)
			}
		}

	}

}
