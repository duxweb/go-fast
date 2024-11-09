package action

import (
	"context"
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

func (t *Resources[T]) Trash(ctx echo.Context) error {
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

	id := ctx.Param("id")
	err = t.trashOne(ctx, id, params)
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Get(ctx, "common.message.trash"),
	})
}

func (t *Resources[T]) TrashBefore(call ActionCallFun[T]) {
	t.trashBeforeFun = call
}

func (t *Resources[T]) TrashAfter(call ActionCallFun[T]) {
	t.trashAfterFun = call
}

func (t *Resources[T]) trashOne(ctx echo.Context, id string, params *gjson.Result) error {
	var model T
	var err error

	err = t.getOne(ctx, &model, id, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BusinessError(i18n.Get(ctx, "common.message.emptyData"))
		} else {
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

	if t.trashBeforeFun != nil {
		err = t.trashBeforeFun(c, &model)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Unscoped().Delete(t.Model, id).Error
	if err != nil {
		return err
	}

	if t.trashAfterFun != nil {
		err = t.trashAfterFun(c, &model)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
