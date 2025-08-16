package action

import (
	"context"
	"reflect"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/resp"
	"gorm.io/gorm/clause"
)

// Create 创建记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Create(ctx context.Context, input *Data) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
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

	// 应用格式化函数
	if res.formatFun != nil {
		formatFn := res.formatFun.(func(*Model, *Data, context.Context) (*Model, error))
		formattedModel, err := formatFn(&model, input, ctx)
		if err != nil {
			return nil, err
		}
		model = *formattedModel
	}

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	c := context.Background()
	c = context.WithValue(c, "tx", tx)

	// 创建前回调
	if res.createBeforeFun != nil {
		createBeforeFn := res.createBeforeFun.(func(context.Context, *Model, *Data) error)
		err := createBeforeFn(c, &model, input)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err := tx.Model(model).Omit(clause.Associations).Create(&model).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建后回调
	if res.createAfterFun != nil {
		createAfterFn := res.createAfterFun.(func(context.Context, *Model, *Data) error)
		err := createAfterFn(c, &model, input)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return resp.Send(ctx, resp.Data[any, resp.EmptyMeta]{
		Message: "创建成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}