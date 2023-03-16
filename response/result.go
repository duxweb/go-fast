package response

import (
	"github.com/gofiber/fiber/v2"
)

type Result struct {
	Ctx *fiber.Ctx
}

func New(ctx *fiber.Ctx) *Result {
	return &Result{Ctx: ctx}
}

type ResultData struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"ok"`
	Data    any    `json:"data"`
}

func (r *Result) Render(name string, bind any) error {
	return r.Ctx.Render(name, bind)
}

func (r *Result) Send(message string, data ...any) error {
	var params any
	if len(data) > 0 {
		params = data[0]
	} else {
		params = map[string]any{}
	}
	res := ResultData{}
	res.Code = 200
	res.Message = message
	res.Data = params
	return r.Ctx.Status(fiber.StatusOK).JSON(res)
}
