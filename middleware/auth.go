package middleware

import (
	"time"

	"github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/config"
	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(app string) echo.MiddlewareFunc {
	key := config.Load("use").GetString("app.secret")
	middle := echojwt.Config{
		SigningKey:  key,
		TokenLookup: "header:Authorization:Bearer ,query:auth",
		ParseTokenFunc: func(c echo.Context, token string) (interface{}, error) {
			data := auth.JwtClaims{}
			jwtToken, err := jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
				return []byte(key), nil
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
			claims := user.Claims.(*auth.JwtClaims)
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

			expire := expiredAt.Sub(issuedAt.Time)

			if expiredAt.Add(-(time.Duration(expire.Seconds()/2) * time.Second)).After(time.Now()) {
				return
			}
			newToken, _ := auth.NewJWT().MakeToken(claims.Subject, claims.ID, expire)
			c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+newToken)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	}
	return echojwt.WithConfig(middle)
}
