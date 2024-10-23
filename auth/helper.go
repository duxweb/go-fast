package auth

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Service struct {
	app string
	ctx *fiber.Ctx
}

func NewService(app string, ctx *fiber.Ctx) *Service {
	return &Service{app: app, ctx: ctx}
}

func (t *Service) ID() string {
	token := t.ctx.Get("Authorization")
	token = strings.ReplaceAll(token, "Bearer ", "")
	if token == "" {
		return ""
	}
	parsingToken, err := NewJWT().ParsingToken(token, t.app)
	if err != nil {
		return ""
	}
	return parsingToken.ID
}
