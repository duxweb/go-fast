package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/resp"
	"github.com/spf13/cast"
)

const apiDiffTime float64 = 10

// ApiMiddleware 签名认证中间件
func ApiMiddleware(secretCallback func(id string) string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// 在开发环境下，如果是文档测试请求，跳过签名验证
		if global.Debug && ctx.Header("X-Doc-Test") == "true" {
			next(ctx)
			return
		}

		// 获取请求头
		date := ctx.Header("Content-Date")
		sign := ctx.Header("Content-MD5")
		id := ctx.Header("AccessKey")

		// 验证时间戳
		timeNow := time.Now()
		t := time.Unix(cast.ToInt64(date), 0)
		if timeNow.Sub(t).Seconds() > apiDiffTime {
			resp.RawSend(ctx, resp.Data[any, any]{
				Code:    408,
				Message: "Request timeout",
			})
			return
		}

		// 获取密钥
		secretKey := secretCallback(id)
		if secretKey == "" {
			resp.RawSend(ctx, resp.Data[any, any]{
				Code:    401,
				Message: "Invalid access key",
			})
			return
		}

		// 构建签名数据
		path := ctx.URL().Path
		queryString := ctx.URL().RawQuery
		signData := []string{
			path,
			queryString,
			date,
		}

		// 计算签名
		h := sha256.New
		mac := hmac.New(h, []byte(secretKey))
		mac.Write([]byte(strings.Join(signData, "\n")))
		digest := mac.Sum(nil)
		hexDigest := hex.EncodeToString(digest)

		// 验证签名
		if sign != hexDigest {
			resp.RawSend(ctx, resp.Data[any, any]{
				Code:    401,
				Message: "Invalid signature",
			})
			return
		}

		// 认证通过，继续处理
		next(ctx)
	}
}
