package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/demdxx/gocast/v2"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

const diffTime float64 = 10

func ApiMiddleware(secretCallback func(id string) string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		date := c.Get("Content-Date")
		timeNow := time.Now()
		t := time.Unix(gocast.Number[int64](date), 0)
		if timeNow.Sub(t).Seconds() > diffTime {
			return fiber.ErrRequestTimeout
		}

		sign := c.Get("Content-MD5")
		id := c.Get("AccessKey")

		secretKey := secretCallback(id)
		if secretKey == "" {
			return fiber.ErrUnauthorized
		}
		signData := []string{
			c.Path(),
			c.Context().QueryArgs().String(),
			date,
		}
		h := sha256.New
		mac := hmac.New(h, []byte(secretKey))
		mac.Write([]byte(strings.Join(signData, "\n")))
		digest := mac.Sum(nil)
		hexDigest := hex.EncodeToString(digest)
		if sign != hexDigest {
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}
}
