package auth

import (
	"errors"
	"github.com/duxweb/go-fast/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"time"
)

func AuthMiddleware(app string, renewals ...time.Duration) echo.MiddlewareFunc {
	var renewal time.Duration = 0
	if len(renewals) > 0 {
		renewal = renewals[0]
	}
	key := config.Load("use").GetString("app.secret")
	middle := echojwt.Config{
		SigningKey:  key,
		TokenLookup: "header:" + echo.HeaderAuthorization + ",query:auth",
		ParseTokenFunc: func(c echo.Context, token string) (interface{}, error) {
			data := JwtClaims{}
			jwtToken, err := jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
				return key, nil
			})
			if err != nil {
				return nil, err
			}
			if data.Subject != app {
				return nil, errors.New("token type error")
			}
			return jwtToken, nil
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(JwtClaims)
			c.Set("auth", claims)
			if !claims.Refresh {
				return
			}

			issuedAt, _ := claims.GetIssuedAt()
			if issuedAt == nil {
				return
			}
			expiredAt, _ := claims.GetExpirationTime()
			if expiredAt == nil {
				return
			}

			if expiredAt.Add(-renewal).After(time.Now()) {
				return
			}
			expire := expiredAt.Sub(issuedAt.Time)
			newToken, _ := NewJWT().MakeToken(claims.Subject, claims.ID, expire)
			c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+newToken)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	}
	return echojwt.WithConfig(middle)
}
