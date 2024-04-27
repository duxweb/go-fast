package auth

import "github.com/labstack/echo/v4"

type Service struct {
	app string
	ctx echo.Context
}

func NewService(app string, ctx echo.Context) *Service {
	return &Service{app: app, ctx: ctx}
}

func (t *Service) ID() string {
	token := t.ctx.Request().Header.Get("Authorization")
	if token == "" {
		return ""
	}
	parsingToken, err := NewJWT().ParsingToken(token, t.app)
	if err != nil {
		return ""
	}
	return parsingToken.ID
}
