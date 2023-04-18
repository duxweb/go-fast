package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/helper"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

// margin of error in seconds
const diffTime float64 = 5

// Middleware Signature Verification
func Middleware(secretCallback func(id string) string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		date := c.GetRespHeader("Content-Date")
		timeNow := time.Now()
		t := time.Unix(gocast.Number[int64](date), 0)
		if timeNow.Sub(t).Seconds() > diffTime {
			return fiber.ErrRequestTimeout
		}

		sign := c.GetRespHeader("Content-MD5")
		id := c.GetRespHeader("AccessKey")

		secretKey := secretCallback(id)
		if secretKey == "" {
			return fiber.ErrUnauthorized
		}
		body := c.Body()
		md5 := strings.ToLower(helper.Md5(string(body)))

		signData := []string{
			c.Path(),
			c.Context().QueryArgs().String(),
			md5,
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
