package action

import (
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
)

func (t *Resources[T]) Create(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(ctx)
		if err != nil {
			return err
		}
	}

	requestData := map[string]any{}
	err = (&echo.DefaultBinder{}).BindBody(ctx, &requestData)
	if err != nil {
		return err
	}

	if t.validatorFun != nil {
		rules, err := t.validatorFun(requestData, ctx)
		if err != nil {
			return err
		}
		err = validator.ValidatorMaps(requestData, rules)
		if err != nil {
			return err
		}
	}

	return response.Send(ctx, response.Data{})
}
