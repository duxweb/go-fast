package action

import (
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (t *Resources[T]) RestoreMany(ctx *fiber.Ctx) error {
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

	ids := strings.Split(params.Get("ids").String(), ",")

	for _, id := range ids {
		if id == "" {
			continue
		}
		err = t.restoreOne(ctx, id)
		if err != nil {
			return err
		}
	}

	return response.Send(ctx, response.Data{
		Message: i18n.Trans.Get("common.message.restore"),
	})
}
