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

// Delete 删除记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Delete(ctx context.Context, input *DeleteInput) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
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

	if res.Tree {
		var count int64
		tx.Model(model).Where("parent_id = ?", id).Count(&count)
		if count > 0 {
			tx.Rollback()
			return nil, errors.New("存在子级数据，不能删除")
		}
	}

	// 删除前回调
	if res.deleteBeforeFun != nil {
		deleteBeforeFn := res.deleteBeforeFun.(func(context.Context, *Model) error)
		err := deleteBeforeFn(c, &model)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Omit(clause.Associations).Delete(model, id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 删除后回调
	if res.deleteAfterFun != nil {
		deleteAfterFn := res.deleteAfterFun.(func(context.Context, *Model) error)
		err := deleteAfterFn(c, &model)
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
		Message: "删除成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}