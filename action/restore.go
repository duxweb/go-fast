package action

import (
	"context"
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (t *Resources[T]) Restore(ctx *fiber.Ctx) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(t, ctx)
		if err != nil {
			return err
		}
	}

	id := ctx.Params("id")
	err = t.restoreOne(ctx, id)
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.restore"),
	})
}

func (t *Resources[T]) RestoreBefore(call ActionCallFun[T]) {
	t.restoreBeforeFun = call
}

func (t *Resources[T]) RestoreAfter(call ActionCallFun[T]) {
	t.restoreAfterFun = call
}

func (t *Resources[T]) restoreOne(ctx *fiber.Ctx, id string) error {
	var model T
	var err error

	err = t.getOne(ctx, &model, id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BusinessError(i18n.Trans.Get("common.message.emptyData"))
		} else {
			return err
		}
	}

	tx := database.Gorm().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	c := context.WithValue(context.Background(), "tx", tx)

	if t.restoreBeforeFun != nil {
		err = t.restoreBeforeFun(c, &model)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Model(&model).Unscoped().Update("deleted_at", nil).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if t.restoreAfterFun != nil {
		err = t.restoreAfterFun(c, &model)
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
