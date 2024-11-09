package action

import (
	"context"
	"errors"

	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (t *Resources[T]) Delete(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(t, ctx)
		if err != nil {
			return err
		}
	}

	id := ctx.Param("id")
	err = t.deleteOne(ctx, id)
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Get(ctx, "common.message.delete"),
	})
}

func (t *Resources[T]) DeleteBefore(call ActionCallFun[T]) {
	t.deleteBeforeFun = call
}

func (t *Resources[T]) DeleteAfter(call ActionCallFun[T]) {
	t.deleteAfterFun = call
}

func (t *Resources[T]) deleteOne(ctx echo.Context, id string) error {
	var model T
	var err error

	err = t.getOne(ctx, &model, id, nil)
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

	if t.Tree {
		var count int64
		tx.Model(t.Model).Where("parent_id = ?", id).Count(&count)
		if count > 0 {
			tx.Rollback()
			return response.BusinessError(i18n.Get(ctx, "common.message.children"))
		}
	}

	if t.deleteBeforeFun != nil {
		err = t.deleteBeforeFun(c, &model)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Delete(t.Model, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if t.deleteAfterFun != nil {
		err = t.deleteAfterFun(c, &model)
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
