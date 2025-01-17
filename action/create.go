package action

import (
	"context"

	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"gorm.io/gorm/clause"
)

func (t *Resources[T]) Create(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(t, ctx)
		if err != nil {
			return err
		}
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
		err = validator.ValidatorMaps(ctx, dataMaps, rules)
		if err != nil {
			return err
		}
	}

	model := t.Model
	if t.formatFun != nil {
		err = t.formatFun(&model, data, ctx)
		if err != nil {
			return err
		}
	}

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	c := context.Background()
	c = context.WithValue(c, "tx", tx)
	c = context.WithValue(c, "echo", ctx)

	if t.createBeforeFun != nil {
		err = t.createBeforeFun(c, &model, data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if t.saveBeforeFun != nil {
		err = t.saveBeforeFun(c, &model, data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Model(t.Model).Omit(clause.Associations).Create(&model).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if t.createAfterFun != nil {
		err = t.createAfterFun(c, &model, data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if t.saveAfterFun != nil {
		err = t.saveAfterFun(c, &model, data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Get(ctx, "common.message.create"),
	})
}

func (t *Resources[T]) CreateBefore(call ActionCallParamsFun[T]) {
	t.createBeforeFun = call
}

func (t *Resources[T]) CreateAfter(call ActionCallParamsFun[T]) {
	t.createAfterFun = call
}

func (t *Resources[T]) SaveBefore(call ActionCallParamsFun[T]) {
	t.saveBeforeFun = call
}

func (t *Resources[T]) SaveAfter(call ActionCallParamsFun[T]) {
	t.saveAfterFun = call
}
