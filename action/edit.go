package action

import (
	"context"
	"reflect"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/resp"
	"gorm.io/gorm/clause"
)

// Edit 编辑记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Edit(ctx context.Context, input *EditInput[Data]) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
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

	// 应用格式化函数
	if res.formatFun != nil {
		formatFn := res.formatFun.(func(*Model, *Data, context.Context) (*Model, error))
		formattedModel, err := formatFn(&model, &input.Body, ctx)
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

	// 编辑前回调
	if res.editBeforeFun != nil {
		editBeforeFn := res.editBeforeFun.(func(context.Context, *Model, *Data) error)
		err := editBeforeFn(c, &model, &input.Body)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Debug().Omit(clause.Associations).Save(&model).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 编辑后回调
	if res.editAfterFun != nil {
		editAfterFn := res.editAfterFun.(func(context.Context, *Model, *Data) error)
		err := editAfterFn(c, &model, &input.Body)
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
		Message: "编辑成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}
