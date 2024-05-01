package action

import (
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/gookit/goutil/structs"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (t *Resources[T]) Store(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(ctx)
		if err != nil {
			return err
		}
	}

	requestData := map[string]any{}
	err = ctx.Bind(&requestData)
	if err != nil {
		return err
	}

	keys := lo.Keys[string, any](requestData)

	if t.validatorFun != nil {
		rules, err := t.validatorFun(requestData, ctx)
		if err != nil {
			return err
		}
		rules = lo.PickBy[string, validator.ValidatorWarp](rules, func(key string, value validator.ValidatorWarp) bool {
			return lo.IndexOf[string](keys, key) != -1
		})
		err = validator.ValidatorMaps(requestData, rules)
		if err != nil {
			return err
		}
	}

	id := ctx.Param("id")
	var model T
	err = t.getOne(ctx, &model, id, requestData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BusinessError(i18n.Trans.Get("common.message.emptyData"))
		} else {
			return err
		}
	}

	if t.formatFun != nil {
		err = t.formatFun(&model, requestData, ctx)
		if err != nil {
			return err
		}
	}

	data, err := structs.StructToMap(model)
	if err != nil {
		return err
	}
	data = lo.PickBy[string, any](data, func(key string, value any) bool {
		return lo.IndexOf[string](keys, key) != -1
	})

	if t.storeBeforeFun != nil {
		err = t.storeBeforeFun(&model, requestData)
		if err != nil {
			return err
		}
	}

	err = database.Gorm().Model(&model).Updates(data).Error
	if err != nil {
		return err
	}

	if t.storeAfterFun != nil {
		err = t.storeAfterFun(&model, requestData)
		if err != nil {
			return err
		}
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.store"),
	})
}

func (t *Resources[T]) StoreBefore(call ActionCallFun[T]) {
	t.storeBeforeFun = call
}

func (t *Resources[T]) StoreAfter(call ActionCallFun[T]) {
	t.storeAfterFun = call
}
