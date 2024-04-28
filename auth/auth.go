package auth

import (
	"errors"
	"github.com/duxweb/go-fast/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type JwtClaims struct {
	Refresh bool `json:"refresh"`
	jwt.RegisteredClaims
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(config.Load("use").GetString("app.secret")),
	}
}

func (j *JWT) MakeToken(app string, id string, expires ...time.Duration) (tokenString string, err error) {
	expire := 86400 * time.Second
	if len(expires) > 0 {
		expire = expires[0]
	}
	claim := JwtClaims{
		Refresh: true,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   app,
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(j.SigningKey)
	return tokenString, err
}

func (j *JWT) ParsingToken(token string, app string) (claims *JwtClaims, err error) {
	data := JwtClaims{}
	_, err = jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if data.Subject != app {
		return nil, errors.New("token type error")
	}
	return &data, nil
}
