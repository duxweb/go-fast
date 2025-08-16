package action

import (
	"context"
	"errors"
	"reflect"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/resp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Trash 软删除（彻底删除）记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Trash(ctx context.Context, input *DeleteInput) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("记录不存在")
		} else {
			return nil, err
		}
	}

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	c := context.Background()
	c = context.WithValue(c, "tx", tx)

	// 彻底删除前回调
	if res.trashBeforeFun != nil {
		trashBeforeFn := res.trashBeforeFun.(func(context.Context, *Model) error)
		err := trashBeforeFn(c, &model)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 使用 Unscoped 进行彻底删除
	err = tx.Unscoped().Omit(clause.Associations).Delete(model, id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 彻底删除后回调
	if res.trashAfterFun != nil {
		trashAfterFn := res.trashAfterFun.(func(context.Context, *Model) error)
		err := trashAfterFn(c, &model)
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
		Message: "彻底删除成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}