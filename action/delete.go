package action

import (
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
	err = t.deleteOne(ctx, id, params)
	if err != nil {
		return err
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.delete"),
	})
}

func (t *Resources[T]) DeleteBefore(call ActionCallFun[T]) {
	t.deleteBeforeFun = call
}

func (t *Resources[T]) DeleteAfter(call ActionCallFun[T]) {
	t.deleteAfterFun = call
}

func (t *Resources[T]) deleteOne(ctx echo.Context, id string, params map[string]any) error {
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

	if t.deleteBeforeFun != nil {
		err = t.deleteBeforeFun(&model, params)
		if err != nil {
			return err
		}
	}

	err = database.Gorm().Delete(t.Model, id).Error
	if err != nil {
		return err
	}

	if t.deleteAfterFun != nil {
		err = t.deleteAfterFun(&model, params)
		if err != nil {
			return err
		}
	}
	return nil
}
