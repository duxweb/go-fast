package action

import (
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
)

func (t *Resources[T]) Create(ctx echo.Context) error {
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

	model := t.Model
	if t.formatFun != nil {
		err = t.formatFun(&model, requestData, ctx)
		if err != nil {
			return err
		}
	}

	if t.createBeforeFun != nil {
		t.createBeforeFun(&model, requestData)
	}
	if t.saveBeforeFun != nil {
		t.saveBeforeFun(&model, requestData)
	}

	err = database.Gorm().Debug().Model(t.Model).Create(&model).Error
	if err != nil {
		return err
	}

	if t.createAfterFun != nil {
		t.createAfterFun(&model, requestData)
	}
	if t.saveAfterFun != nil {
		t.saveAfterFun(&model, requestData)
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.create"),
	})
}

func (t *Resources[T]) CreateBefore(call ActionCallFun[T]) {
	t.createBeforeFun = call
}

func (t *Resources[T]) CreateAfter(call ActionCallFun[T]) {
	t.createAfterFun = call
}

func (t *Resources[T]) SaveBefore(call ActionCallFun[T]) {
	t.saveBeforeFun = call
}

func (t *Resources[T]) SaveAfter(call ActionCallFun[T]) {
	t.saveAfterFun = call
}