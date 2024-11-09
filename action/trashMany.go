package action

import (
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"strings"
)

func (t *Resources[T]) TrashMany(ctx echo.Context) error {
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

	ids := strings.Split(data.Get("ids").String(), ",")

	for _, id := range ids {
		if id == "" {
			continue
		}
		err = t.trashOne(ctx, id, params)
		if err != nil {
			return err
		}
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Get(ctx, "common.message.trash"),
	})
}
