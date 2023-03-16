package request

import (
	"github.com/duxweb/go-fast/validator"
	"github.com/gofiber/fiber/v2"
)

func BodyParser(ctx *fiber.Ctx, params any) error {
	var err error
	if err = ctx.BodyParser(params); err != nil {
		return err
	}
	err = validator.Validator().Struct(params)
	if err = validator.ProcessError(params, err); err != nil {
		return err
	}
	return nil
}

func QueryParser(ctx *fiber.Ctx, params any) error {
	var err error
	if err = ctx.QueryParser(params); err != nil {
		return err
	}
	err = validator.Validator().Struct(params)
	if err = validator.ProcessError(params, err); err != nil {
		return err
	}
	return nil
}
