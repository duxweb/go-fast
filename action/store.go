package action

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/resp"
	"github.com/gookit/goutil/structs"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Store 部分字段更新方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) Store(ctx context.Context, input *EditInput[Data]) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
	id := input.ID
	
	// 将输入转换为 JSON 以获取实际提交的字段
	inputBytes, _ := json.Marshal(input.Body)
	var inputMap map[string]interface{}
	json.Unmarshal(inputBytes, &inputMap)
	
	// 获取提交的字段名
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}

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

	// 应用格式化函数
	if res.formatFun != nil {
		formatFn := res.formatFun.(func(*Model, *Data, context.Context) (*Model, error))
		formattedModel, err := formatFn(&model, &input.Body, ctx)
		if err != nil {
			return nil, err
		}
		model = *formattedModel
	}

	// 将模型转换为 map，然后只保留提交的字段
	formatData, err := structs.StructToMap(model)
	if err != nil {
		return nil, err
	}
	formatData = lo.PickBy[string, any](formatData, func(key string, value any) bool {
		return lo.IndexOf[string](keys, key) != -1
	})

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	c := context.Background()
	c = context.WithValue(c, "tx", tx)

	// Store前回调
	if res.storeBeforeFun != nil {
		storeBeforeFn := res.storeBeforeFun.(func(context.Context, *Model, *Data) error)
		err := storeBeforeFn(c, &model, &input.Body)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Model(&model).Omit(clause.Associations).Updates(formatData).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Store后回调
	if res.storeAfterFun != nil {
		storeAfterFn := res.storeAfterFun.(func(context.Context, *Model, *Data) error)
		err := storeAfterFn(c, &model, &input.Body)
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
		Message: "保存成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}