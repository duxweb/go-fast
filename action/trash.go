package action

import (
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (t *Resources[T]) Trash(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(ctx)
		if err != nil {
			return err
		}
	}

	params := map[string]any{}
	err = ctx.Bind(&params)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	err = t.trashOne(ctx, id, params)
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.trash"),
	})
}

func (t *Resources[T]) TrashBefore(call ActionCallFun[T]) {
	t.trashBeforeFun = call
}

func (t *Resources[T]) TrashAfter(call ActionCallFun[T]) {
	t.trashAfterFun = call
}

func (t *Resources[T]) trashOne(ctx echo.Context, id string, params map[string]any) error {
	var model T
	var err error

	err = t.getOne(ctx, &model, id, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BusinessError(i18n.Trans.Get("common.message.emptyData"))
		} else {
			return err
		}
	}

	if t.trashBeforeFun != nil {
		t.trashBeforeFun(&model, params)
	}

	err = database.Gorm().Unscoped().Delete(t.Model, id).Error
	if err != nil {
		return err
	}

	if t.trashAfterFun != nil {
		t.trashAfterFun(&model, params)
	}
	return nil
}
