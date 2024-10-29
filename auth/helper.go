package auth

import (
	"github.com/labstack/echo/v4"
	"strings"
)

type Service struct {
	app string
	ctx echo.Context
}

func NewService(app string, ctx echo.Context) *Service {
	return &Service{app: app, ctx: ctx}
}

func (t *Service) ID() string {
	token := t.ctx.Request().Header.Get("Authorization")
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
