package action

import (
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"strings"
)

func (t *Resources[T]) TrashMany(ctx echo.Context) error {
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

	ids := strings.Split(cast.ToString(params["ids"]), ",")

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
		Message: i18n.Trans.Get("common.message.trash"),
	})
}
