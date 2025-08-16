package middleware

import (
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/auth"
	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

func AuthMiddleware(app string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// 获取 token
		token := ctx.Header("Authorization")
		if token == "" {
			token = ctx.Query("auth")
		} else {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			ctx.SetStatus(401)
			return
		}

		// 解析 token
		key := config.Load("use").String("app.secret")
		data := auth.JwtClaims{}
		jwtToken, err := jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		if err != nil || !jwtToken.Valid {
			ctx.SetStatus(401)
			return
		}

		if data.Subject != app {
			ctx.SetStatus(401)
			return
		}

		// 将认证信息存入上下文
		ctx = huma.WithValue(ctx, "auth", &data)
		ctx = huma.WithValue(ctx, "user", jwtToken)

		// Token 自动刷新逻辑
		if data.Refresh {
			issuedAt, _ := data.GetIssuedAt()
			expiredAt, _ := data.GetExpirationTime()

			if issuedAt != nil && expiredAt != nil {
				expire := expiredAt.Sub(issuedAt.Time)

				// 如果 token 过期时间过半，则刷新 token
				if time.Now().After(expiredAt.Add(-(time.Duration(expire.Seconds()/2) * time.Second))) {
					newToken, _ := auth.NewJWT().MakeToken(data.Subject, data.ID, expire)
					ctx.SetHeader("Authorization", "Bearer "+newToken)
				}
			}
		}

		next(ctx)
	}
}

func Can(ctx huma.Context, name string) error {
	userPermission, _ := ctx.Context().Value("userPermissions").([]string)
	permission, _ := ctx.Context().Value("permissions").([]string)

	if len(userPermission) == 0 || len(permission) == 0 {
		return nil
	}

	if lo.IndexOf[string](permission, name) == -1 {
		return nil
	}

	if lo.IndexOf[string](userPermission, name) != -1 {
		return nil
	}

	return errors.NewHTTPError(403, "Forbidden")
}
