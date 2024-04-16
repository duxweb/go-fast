package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/demdxx/gocast/v2"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

const diffTime float64 = 10

func ApiMiddleware(secretCallback func(id string) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			date := c.Response().Header().Get("Content-Date")
			timeNow := time.Now()
			t := time.Unix(gocast.Number[int64](date), 0)
			if timeNow.Sub(t).Seconds() > diffTime {
				return echo.ErrRequestTimeout
			}

			sign := c.Request().Header.Get("Content-MD5")
			id := c.Request().Header.Get("AccessKey")

			secretKey := secretCallback(id)
			if secretKey == "" {
				return echo.ErrUnauthorized
			}
			signData := []string{
				c.Path(),
				c.QueryString(),
				date,
			}
			h := sha256.New
			mac := hmac.New(h, []byte(secretKey))
			mac.Write([]byte(strings.Join(signData, "\n")))
			digest := mac.Sum(nil)
			hexDigest := hex.EncodeToString(digest)
			if sign != hexDigest {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
