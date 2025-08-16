package middleware

import (
	"github.com/danielgtaylor/huma/v2"
)

// OperateMiddleware 操作日志中间件
func OperateMiddleware(name string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// TODO: 这里可以记录操作日志
		// 例如：记录用户访问了哪个API，请求参数等
		
		// 继续执行下一个中间件
		next(ctx)
	}
}