package action

import (
	"context"
	"errors"
	"reflect"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/resp"
	"gorm.io/gorm"
)

// Restore 恢复软删除记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Restore(ctx context.Context, input *DeleteInput) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
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

	// 恢复前回调
	if res.restoreBeforeFun != nil {
		restoreBeforeFn := res.restoreBeforeFun.(func(context.Context, *Model) error)
		err := restoreBeforeFn(c, &model)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 恢复软删除的数据
	err = tx.Model(&model).Unscoped().Update("deleted_at", nil).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 恢复后回调
	if res.restoreAfterFun != nil {
		restoreAfterFn := res.restoreAfterFun.(func(context.Context, *Model) error)
		err := restoreAfterFn(c, &model)
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
		Message: "恢复成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}