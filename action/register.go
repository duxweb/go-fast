package action

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/route"
)

// RegisterRoutes 自动注册所有 CRUD 路由
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) RegisterRoutes(appName, resName, routePath string) {
	// 获取路由组
	routeData := route.GetRouter(appName)
	fmt.Println("xx", routeData)
	if routeData == nil {
		return // 如果没有找到路由组，直接返回
	}

	group := routeData

	// 注册列表路由
	if res.ActionList {
		route.Get(group, routePath, "list", huma.Operation{
			OperationID: resName + ".list",
			Summary:     "获取" + resName + "列表",
			Tags:        []string{resName},
		}, res.List)
	}

	// 注册详情路由
	if res.ActionShow {
		route.Get(group, routePath+"/{id}", "show", huma.Operation{
			OperationID: resName + ".show",
			Summary:     "获取" + resName + "详情",
			Tags:        []string{resName},
		}, res.Show)
	}

	// 注册创建路由
	if res.ActionCreate {
		route.Post(group, routePath, "create", huma.Operation{
			OperationID: resName + ".create",
			Summary:     "创建" + resName,
			Tags:        []string{resName},
		}, res.Create)
	}

	// 注册编辑路由
	if res.ActionEdit {
		route.Put(group, routePath+"/{id}", "edit", huma.Operation{
			OperationID: resName + ".edit",
			Summary:     "编辑" + resName,
			Tags:        []string{resName},
		}, res.Edit)
	}

	// 注册单字段更新路由
	if res.ActionStore {
		route.Patch(group, routePath+"/{id}", "store", huma.Operation{
			OperationID: resName + ".store",
			Summary:     "更新" + resName + "字段",
			Tags:        []string{resName},
		}, res.Store)
	}

	// 注册删除路由
	if res.ActionDelete {
		route.Delete(group, routePath+"/{id}", "delete", huma.Operation{
			OperationID: resName + ".delete",
			Summary:     "删除" + resName,
			Tags:        []string{resName},
		}, res.Delete)

		route.Post(group, routePath+"/delete", "deleteMany", huma.Operation{
			OperationID: resName + ".deleteMany",
			Summary:     "批量删除" + resName,
			Tags:        []string{resName},
		}, res.DeleteMany)
	}

	// 注册彻底删除路由
	if res.ActionTrash {
		route.Delete(group, routePath+"/{id}/trash", "trash", huma.Operation{
			OperationID: resName + ".trash",
			Summary:     "彻底删除" + resName,
			Tags:        []string{resName},
		}, res.Trash)

		route.Post(group, routePath+"/trash", "trashMany", huma.Operation{
			OperationID: resName + ".trashMany",
			Summary:     "批量彻底删除" + resName,
			Tags:        []string{resName},
		}, res.TrashMany)
	}

	// 注册恢复路由
	if res.ActionRestore {
		route.Post(group, routePath+"/{id}/restore", "restore", huma.Operation{
			OperationID: resName + ".restore",
			Summary:     "恢复" + resName,
			Tags:        []string{resName},
		}, res.Restore)

		route.Post(group, routePath+"/restore", "restoreMany", huma.Operation{
			OperationID: resName + ".restoreMany",
			Summary:     "批量恢复" + resName,
			Tags:        []string{resName},
		}, res.RestoreMany)
	}
}
