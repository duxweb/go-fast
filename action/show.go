package action

import (
	"context"
	"reflect"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/resp"
	"gorm.io/gorm"
)

// Show 详情查询方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Show(ctx context.Context, input *ShowInput) (*resp.HumaResponse[Info, DetailMeta], error) {
	id := input.ID

	// 获取模型实例
	var model Model
	if res.model != nil {
		// 使用反射创建模型实例
		modelType := reflect.TypeOf(res.model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		model = reflect.New(modelType).Interface().(Model)
	}

	err := res.getOne(&model, id, ctx)
	if err != nil {
		return nil, err
	}

	var data Info
	if res.transformFun != nil {
		transformFn := res.transformFun.(func(*Model, int, *TransformContext) Info)
		data = transformFn(&model, 0, &TransformContext{IsList: false})
	}

	// 处理meta数据
	var metaStruct DetailMeta
	if res.metaOneFun != nil {
		if metaFn, ok := res.metaOneFun.(func(Model) DetailMeta); ok {
			metaStruct = metaFn(model)
		}
	} else {
		// 使用默认空meta，需要转换为DetailMeta
		if m, ok := any(resp.EmptyMeta{}).(DetailMeta); ok {
			metaStruct = m
		}
	}

	return resp.Send(ctx, resp.Data[Info, DetailMeta]{
		Data: data,
		Meta: metaStruct,
	}), nil
}

// getOne 获取单条记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) getOne(model *Model, id string, ctx context.Context) error {
	query := database.Gorm().Unscoped().Model(model).Where(res.Key+" = ?", id)

	// 应用全局查询
	if res.queryFun != nil {
		queryFn := res.queryFun.(func(*gorm.DB, context.Context) *gorm.DB)
		query = queryFn(query, ctx)
	}

	if res.preload != nil {
		for _, v := range res.preload {
			query = query.Preload(v)
		}
	}
	return query.First(model).Error
}
