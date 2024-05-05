package action

import (
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (t *Resources[T]) Restore(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(t, ctx)
		if err != nil {
			return err
		}
	}

	id := ctx.Param("id")
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

func (t *Resources[T]) restoreOne(ctx echo.Context, id string) error {
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

	if t.restoreBeforeFun != nil {
		err = t.restoreBeforeFun(&model)
		if err != nil {
			return err
		}
	}

	err = database.Gorm().Model(&model).Unscoped().Update("deleted_at", nil).Error
	if err != nil {
		return err
	}

	if t.restoreAfterFun != nil {
		err = t.restoreAfterFun(&model)
		if err != nil {
			return err
		}
	}
	return nil
}
