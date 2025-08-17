package action

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/route"
)

// Register 自动注册所有 CRUD 路由
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Register(appName, resName, routePath string) *route.RouterData {
	// 获取路由组
	routeData := route.GetRouter(appName)

	group := route.Group(routeData, routePath, resName, res.routeLabel)

	// 注册列表路由
	if res.ActionList {
		route.Get(group, "", "list", huma.Operation{
			Summary: "列表",
		}, res.List)
	}

	// 注册详情路由
	if res.ActionShow {
		route.Get(group, "/{id}", "show", huma.Operation{
			Summary: "详情",
		}, res.Show)
	}

	// 注册创建路由
	if res.ActionCreate {
		route.Post(group, "", "create", huma.Operation{
			Summary: "创建",
		}, res.Create)
	}

	// 注册编辑路由
	if res.ActionEdit {
		route.Put(group, "/{id}", "edit", huma.Operation{
			Summary: "编辑",
		}, res.Edit)
	}

	// 注册单字段更新路由
	if res.ActionStore {
		route.Patch(group, "/{id}", "store", huma.Operation{
			Summary: "更新",
		}, res.Store)
	}

	// 注册删除路由
	if res.ActionDelete {
		route.Delete(group, "/{id}", "delete", huma.Operation{
			Summary: "删除",
		}, res.Delete)

		route.Delete(group, "", "deleteMany", huma.Operation{
			Summary: "批量删除",
		}, res.DeleteMany)
	}

	// 注册彻底删除路由
	if res.ActionTrash {
		route.Delete(group, "/{id}/trash", "trash", huma.Operation{
			Summary: "清空",
		}, res.Trash)

		route.Delete(group, "/trash", "trashMany", huma.Operation{
			Summary: "批量清空",
		}, res.TrashMany)
	}

	// 注册恢复路由
	if res.ActionRestore {
		route.Post(group, "/{id}/restore", "restore", huma.Operation{
			Summary: "恢复",
		}, res.Restore)

		route.Post(group, "/restore", "restoreMany", huma.Operation{
			Summary: "批量恢复",
		}, res.RestoreMany)
	}

	return group
}
