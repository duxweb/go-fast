package action

// Result action操作函数映射类型
type Result = map[string]any

// TransformContext Transform回调的上下文
type TransformContext struct {
	IsList bool `json:"is_list"` // 是否为列表操作
	// 可扩展其他字段，如 IsExport bool、FilterFields []string 等
}

// Meta 列表元数据结构
type Meta struct {
	Total int `json:"total,omitempty" doc:"总数"`
	Page  int `json:"page,omitempty" doc:"当前页"`
	Limit int `json:"limit,omitempty" doc:"每页数量"`
}

type Pagination struct {
	Status   bool
	PageSize int
}

// Resources 资源配置结构体（泛型版本）
type Resources[Model any, Info any, Params any, Data any, ListMeta any, DetailMeta any] struct {
	Key        string
	Tree       bool
	TreeSort   string
	Pagination Pagination
	preload    []string

	// 路由注册配置
	routeAppName string
	routeResName string
	routePath    string

	// 存储各种回调函数，使用具体类型
	model        any
	queryFun     any // func(*gorm.DB, context.Context) *gorm.DB
	filterFun    any // func(*gorm.DB, *Params) *gorm.DB
	transformFun any // func(*Model, int, *TransformContext) Info
	formatFun    any // func(*Model, *Data, context.Context) (*Model, error)
	metaManyFun  any // func([]Model) ListMeta 或 func([]Model) map[string]any
	metaOneFun   any // func(Model) DetailMeta 或 func(Model) map[string]any

	// 生命周期回调
	createBeforeFun  any
	createAfterFun   any
	editBeforeFun    any
	editAfterFun     any
	saveBeforeFun    any
	saveAfterFun     any
	storeBeforeFun   any
	storeAfterFun    any
	deleteBeforeFun  any
	deleteAfterFun   any
	trashBeforeFun   any
	trashAfterFun    any
	restoreBeforeFun any
	restoreAfterFun  any

	// 操作开关配置项
	ActionList       bool
	ActionShow       bool
	ActionCreate     bool
	ActionEdit       bool
	ActionDelete     bool
	ActionStore      bool
	ActionTrash      bool
	ActionRestore    bool
	ActionSoftDelete bool
	Extend           map[string]any
}

// New 创建资源实例（泛型版本）
func New[Model any, Info any, Params any, Data any, ListMeta any, DetailMeta any]() *Resources[Model, Info, Params, Data, ListMeta, DetailMeta] {
	return &Resources[Model, Info, Params, Data, ListMeta, DetailMeta]{
		Key:  "id",
		Tree: false,
		Pagination: Pagination{
			Status:   true,
			PageSize: 10,
		},
		ActionList:       true,
		ActionShow:       true,
		ActionCreate:     true,
		ActionEdit:       true,
		ActionDelete:     true,
		ActionStore:      true,
		ActionTrash:      false,
		ActionRestore:    false,
		ActionSoftDelete: false,
		Extend:           map[string]any{},
	}
}

// ShowInput 详情查询输入
type ShowInput struct {
	ID string `path:"id" required:"true" doc:"记录ID"`
}

// EditInput 编辑操作输入
type EditInput[Data any] struct {
	ID   string `path:"id" required:"true" doc:"记录ID"`
	Body Data   `json:",inline"`
}

// DeleteInput 删除操作输入
type DeleteInput struct {
	ID string `path:"id" required:"true" doc:"记录ID"`
}

// DeleteManyInput 批量删除输入
type DeleteManyInput struct {
	IDs []string `json:"ids" required:"true" doc:"记录ID列表"`
}

// TrashManyInput 批量彻底删除输入
type TrashManyInput struct {
	IDs []string `json:"ids" required:"true" doc:"记录ID列表"`
}

// RestoreManyInput 批量恢复输入
type RestoreManyInput struct {
	IDs []string `json:"ids" required:"true" doc:"记录ID列表"`
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetModel() any {
	return res.model
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetFormatFun() any {
	return res.formatFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetFilterFun() any {
	return res.filterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetTransformFun() any {
	return res.transformFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetMetaManyFun() any {
	return res.metaManyFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetMetaOneFun() any {
	return res.metaOneFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetCreateBeforeFun() any {
	return res.createBeforeFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetCreateAfterFun() any {
	return res.createAfterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetEditBeforeFun() any {
	return res.editBeforeFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetEditAfterFun() any {
	return res.editAfterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetDeleteBeforeFun() any {
	return res.deleteBeforeFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetDeleteAfterFun() any {
	return res.deleteAfterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetStoreBeforeFun() any {
	return res.storeBeforeFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetStoreAfterFun() any {
	return res.storeAfterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetTrashBeforeFun() any {
	return res.trashBeforeFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetTrashAfterFun() any {
	return res.trashAfterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetRestoreBeforeFun() any {
	return res.restoreBeforeFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetRestoreAfterFun() any {
	return res.restoreAfterFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) IsTree() bool {
	return res.Tree
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetTreeSort() string {
	return res.TreeSort
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetPagination() Pagination {
	return res.Pagination
}

// Getter 方法用于接口访问
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetKey() string {
	return res.Key
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetQueryFun() any {
	return res.queryFun
}

func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) GetPreload() []string {
	return res.preload
}

// Result 返回所有启用的操作函数
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Result() map[string]any {
	result := make(map[string]any)

	if res.ActionList {
		result["list"] = res.List
	}
	if res.ActionShow {
		result["show"] = res.Show
	}
	if res.ActionCreate {
		result["create"] = res.Create
	}
	if res.ActionEdit {
		result["edit"] = res.Edit
	}
	if res.ActionStore {
		result["store"] = res.Store
	}
	if res.ActionDelete {
		result["delete"] = res.Delete
		result["deleteMany"] = res.DeleteMany
	}
	if res.ActionTrash {
		result["trash"] = res.Trash
		result["trashMany"] = res.TrashMany
	}
	if res.ActionRestore {
		result["restore"] = res.Restore
		result["restoreMany"] = res.RestoreMany
	}

	return result
}

// SetRoute 设置路由注册配置
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) SetRoute(appName, resName, routePath string) *Resources[Model, Info, Params, Data, ListMeta, DetailMeta] {
	res.routeAppName = appName
	res.routeResName = resName
	res.routePath = routePath
	return res
}

// Register 执行路由注册
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Register() {
	if res.routeAppName != "" && res.routeResName != "" && res.routePath != "" {
		res.RegisterRoutes(res.routeAppName, res.routeResName, res.routePath)
	}
}
