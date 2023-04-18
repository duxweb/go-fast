package auth

import (
	"errors"
	"github.com/duxweb/go-fast/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWT struct {
	SigningKey []byte
}

// NewJWT Authorization Generation and Decoding
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(config.Get("app").GetString("app.safeKey")),
	}
}

func (j *JWT) MakeToken(app string, params jwt.MapClaims, expires ...int64) (tokenString string, err error) {
	var expire int64 = 86400
	if len(expires) > 0 {
		expire = expires[0]
	}
	claim := jwt.MapClaims{
		"sub": app,
		"exp": time.Now().Add(time.Duration(expire) * time.Minute).Unix(), // Expiration Time
		"iat": time.Now().Unix(),                                          // Issued At Time
	}
	for key, value := range params {
		claim[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(j.SigningKey)
	return tokenString, err
}

func (j *JWT) ParsingToken(token string, app ...string) (claims jwt.MapClaims, err error) {
	data := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if len(app) > 0 && data["sub"] != app[0] {
		return nil, errors.New("token type error")
	}
	return data, nil
}
