package middleware

import (
	"github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/config"
	"github.com/go-errors/errors"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func AuthMiddleware(app string, renewals ...time.Duration) fiber.Handler {
	var renewal time.Duration = 3600
	if len(renewals) > 0 {
		renewal = renewals[0]
	}
	key := config.Load("use").GetString("app.secret")
	middle := jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(key)},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(*auth.JwtClaims)

			if claims.Subject != app {
				return errors.New("token type error")
			}
			c.Locals("auth", claims)
			if !claims.Refresh {
				return nil
			}

			issuedAt, _ := claims.GetIssuedAt()
			if issuedAt == nil {
				return nil
			}
			expiredAt, _ := claims.GetExpirationTime()
			if expiredAt == nil {
				return nil
			}

			if expiredAt.Add(-renewal).After(time.Now()) {
				return nil
			}
			expire := expiredAt.Sub(issuedAt.Time)
			newToken, _ := auth.NewJWT().MakeToken(claims.Subject, claims.ID, expire)
			c.Set(fiber.HeaderAuthorization, "Bearer "+newToken)
			return nil
		},
	}

	return jwtware.New(middle)
}
