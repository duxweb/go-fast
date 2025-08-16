package action

import (
	"context"

	"gorm.io/gorm"
)

// SetModel 设置模型
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) SetModel(model Model) {
	res.model = model
}

// Query 设置全局查询
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Query(fn func(tx *gorm.DB, ctx context.Context) *gorm.DB) {
	res.queryFun = fn
}

// Filter 设置筛选查询（带参数）
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Filter(fn func(tx *gorm.DB, params *Params) *gorm.DB) {
	res.filterFun = fn
}

// Transform 设置数据转换
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Transform(fn func(item *Model, index int, ctx *TransformContext) Info) {
	res.transformFun = fn
}

// Format 设置数据格式化
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Format(fn func(model *Model, data *Data, ctx context.Context) (*Model, error)) {
	res.formatFun = fn
}

// MetaMany 设置列表元数据
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) MetaMany(fn func(data []Model) map[string]any) {
	res.metaManyFun = fn
}

// ListMeta 设置列表元数据（支持自定义Meta类型）
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) ListMeta(fn func(data []Model) ListMeta) {
	res.metaManyFun = fn
}

// MetaOne 设置详情元数据
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) MetaOne(fn func(data Model) DetailMeta) {
	res.metaOneFun = fn
}

// Preload 设置预加载关系
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Preload(relations ...string) {
	res.preload = relations
}

// CreateBefore 创建前回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) CreateBefore(fn func(ctx context.Context, model *Model, data *Data) error) {
	res.createBeforeFun = fn
}

// CreateAfter 创建后回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) CreateAfter(fn func(ctx context.Context, model *Model, data *Data) error) {
	res.createAfterFun = fn
}

// EditBefore 编辑前回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) EditBefore(fn func(ctx context.Context, model *Model, data *Data) error) {
	res.editBeforeFun = fn
}

// EditAfter 编辑后回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) EditAfter(fn func(ctx context.Context, model *Model, data *Data) error) {
	res.editAfterFun = fn
}

// DeleteBefore 删除前回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) DeleteBefore(fn func(ctx context.Context, model *Model) error) {
	res.deleteBeforeFun = fn
}

// DeleteAfter 删除后回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) DeleteAfter(fn func(ctx context.Context, model *Model) error) {
	res.deleteAfterFun = fn
}

// StoreBefore Store前回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) StoreBefore(fn func(ctx context.Context, model *Model, data *Data) error) {
	res.storeBeforeFun = fn
}

// StoreAfter Store后回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) StoreAfter(fn func(ctx context.Context, model *Model, data *Data) error) {
	res.storeAfterFun = fn
}

// TrashBefore 彻底删除前回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) TrashBefore(fn func(ctx context.Context, model *Model) error) {
	res.trashBeforeFun = fn
}

// TrashAfter 彻底删除后回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) TrashAfter(fn func(ctx context.Context, model *Model) error) {
	res.trashAfterFun = fn
}

// RestoreBefore 恢复前回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) RestoreBefore(fn func(ctx context.Context, model *Model) error) {
	res.restoreBeforeFun = fn
}

// RestoreAfter 恢复后回调
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) RestoreAfter(fn func(ctx context.Context, model *Model) error) {
	res.restoreAfterFun = fn
}
