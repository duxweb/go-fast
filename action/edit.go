package action

import (
	"context"
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/models"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

func (t *Resources[T]) Edit(ctx *fiber.Ctx) error {
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

	id := ctx.Params("id")
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

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	c := context.WithValue(context.Background(), "tx", tx)

	if t.editBeforeFun != nil {
		err = t.editBeforeFun(c, &model, data)
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

	if t.Tree && data.Get("parent_id").Exists() {
		parentId := data.Get("parent_id").Uint()
		if !models.CheckParentHas[T](database.Gorm().Model(t.Model), cast.ToUint(id), uint(parentId)) {
			tx.Rollback()
			return response.BusinessError(i18n.Trans.Get("common.message.parent"))
		}
	}

	err = tx.Save(&model).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if t.editAfterFun != nil {
		err = t.editAfterFun(c, &model, data)
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
		Message: i18n.Trans.Get("common.message.edit"),
	})
}

func (t *Resources[T]) EditBefore(call ActionCallParamsFun[T]) {
	t.editBeforeFun = call
}

func (t *Resources[T]) EditAfter(call ActionCallParamsFun[T]) {
	t.editAfterFun = call
}
