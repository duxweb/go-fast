package action

import (
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
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

	requestData := map[string]any{}
	err = ctx.Bind(&requestData)
	if err != nil {
		return err
	}

	if t.validatorFun != nil {
		rules, err := t.validatorFun(requestData, ctx)
		if err != nil {
			return err
		}
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

	if t.editBeforeFun != nil {
		err = t.editBeforeFun(&model, requestData)
		if err != nil {
			return err
		}
	}
	if t.saveBeforeFun != nil {
		err = t.saveBeforeFun(&model, requestData)
		if err != nil {
			return err
		}
	}

	err = database.Gorm().Save(&model).Error
	if err != nil {
		return err
	}

	if t.editAfterFun != nil {
		err = t.editAfterFun(&model, requestData)
		if err != nil {
			return err
		}
	}
	if t.saveAfterFun != nil {
		err = t.saveAfterFun(&model, requestData)
		if err != nil {
			return err
		}
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.edit"),
	})
}

func (t *Resources[T]) EditBefore(call ActionCallFun[T]) {
	t.editBeforeFun = call
}

func (t *Resources[T]) EditAfter(call ActionCallFun[T]) {
	t.editAfterFun = call
}
