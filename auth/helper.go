package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Service 泛型服务
type Service[T any] struct {
	app string
	ctx T
}

func NewService[T any](app string, ctx T) *Service[T] {
	return &Service[T]{app: app, ctx: ctx}
}

func (t *Service[T]) ID() string {
	var token string

	// 根据 context 类型来获取 Authorization header 和 cookie
	switch c := any(t.ctx).(type) {
	case echo.Context:
		// 优先从 Authorization header 获取
		token = c.Request().Header.Get("Authorization")
		// 如果 header 中没有，尝试从 cookie 获取
		if token == "" {
			if cookie, err := c.Cookie("authorization"); err == nil {
				token = "Bearer " + cookie.Value
			}
		}
	case context.Context:
		// 尝试从 Huma context 中获取 http.Request
		var req *http.Request
		if r, ok := c.Value("request").(*http.Request); ok {
			req = r
		} else if r, ok := c.Value(http.ServerContextKey).(*http.Request); ok {
			req = r
		}

		if req != nil {
			// 优先从 Authorization header 获取
			token = req.Header.Get("Authorization")
			// 如果 header 中没有，尝试从 cookie 获取
			if token == "" {
				if cookie, err := req.Cookie("authorization"); err == nil {
					token = "Bearer " + cookie.Value
				}
			}
		}
	default:
		return ""
	}

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

// GetID 简化的泛型函数
func GetID[T any](ctx T, app string) string {
	return NewService(app, ctx).ID()
}
