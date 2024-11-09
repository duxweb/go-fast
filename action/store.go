package action

import (
	"context"
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/gookit/goutil/structs"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

func (t *Resources[T]) Store(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(t, ctx)
		if err != nil {
			return err
		}
	}

	params, err := helper.Qs(ctx)
	if err != nil {
		return err
	}

	data, err := helper.Body(ctx)
	if err != nil {
		return err
	}

	keys := []string{}
	data.ForEach(func(key, value gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})

	if t.validatorFun != nil {
		rules, err := t.validatorFun(data, ctx)
		if err != nil {
			return err
		}
		rules = lo.PickBy[string, validator.ValidatorWarp](rules, func(key string, value validator.ValidatorWarp) bool {
			return lo.IndexOf[string](keys, key) != -1
		})
		dataMaps := map[string]any{}
		data.ForEach(func(key, value gjson.Result) bool {
			dataMaps[key.String()] = value.Value()
			return true
		})
		err = validator.ValidatorMaps(ctx, dataMaps, rules)
		if err != nil {
			return err
		}
	}

	id := ctx.Param("id")
	var model T
	err = t.getOne(ctx, &model, id, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BusinessLangError(ctx, "common.message.emptyData")
		} else {
			return err
		}
	}

	if t.formatFun != nil {
		err = t.formatFun(&model, data, ctx)
		if err != nil {
			return err
		}
	}

	formatData, err := structs.StructToMap(model)
	if err != nil {
		return err
	}
	formatData = lo.PickBy[string, any](formatData, func(key string, value any) bool {
		return lo.IndexOf[string](keys, key) != -1
	})

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	c := context.Background()
	c = context.WithValue(c, "tx", tx)
	c = context.WithValue(c, "echo", ctx)

	if t.storeBeforeFun != nil {
		err = t.storeBeforeFun(c, &model, data)
		if err != nil {
			return err
		}
	}

	err = tx.Model(&model).Updates(formatData).Error
	if err != nil {
		return err
	}

	if t.storeAfterFun != nil {
		err = t.storeAfterFun(c, &model, data)
		if err != nil {
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Get(ctx, "common.message.store"),
	})
}

func (t *Resources[T]) StoreBefore(call ActionCallParamsFun[T]) {
	t.storeBeforeFun = call
}

func (t *Resources[T]) StoreAfter(call ActionCallParamsFun[T]) {
	t.storeAfterFun = call
}
