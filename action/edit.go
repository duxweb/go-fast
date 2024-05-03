package action

import (
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

func (t *Resources[T]) Edit(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(ctx)
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

	if t.validatorFun != nil {
		rules, err := t.validatorFun(data, ctx)
		if err != nil {
			return err
		}
		dataMaps := map[string]any{}
		data.ForEach(func(key, value gjson.Result) bool {
			dataMaps[key.String()] = value.Value()
			return true
		})
		err = validator.ValidatorMaps(dataMaps, rules)
		if err != nil {
			return err
		}
	}

	id := ctx.Param("id")
	var model T

	err = t.getOne(ctx, &model, id, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BusinessError(i18n.Trans.Get("common.message.emptyData"))
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

	if t.editBeforeFun != nil {
		err = t.editBeforeFun(&model, data)
		if err != nil {
			return err
		}
	}
	if t.saveBeforeFun != nil {
		err = t.saveBeforeFun(&model, data)
		if err != nil {
			return err
		}
	}

	err = database.Gorm().Save(&model).Error
	if err != nil {
		return err
	}

	if t.editAfterFun != nil {
		err = t.editAfterFun(&model, data)
		if err != nil {
			return err
		}
	}
	if t.saveAfterFun != nil {
		err = t.saveAfterFun(&model, data)
		if err != nil {
			return err
		}
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.edit"),
	})
}

func (t *Resources[T]) EditBefore(call ActionCallParamsFun[T]) {
	t.editBeforeFun = call
}

func (t *Resources[T]) EditAfter(call ActionCallParamsFun[T]) {
	t.editAfterFun = call
}
