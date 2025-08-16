package route

import (
	"github.com/danielgtaylor/huma/v2"
)

func GetRouteName(ctx huma.Context) string {
	// 在 Huma 中，可以通过 Operation 获取路由信息
	if ctx.Operation() != nil {
		return ctx.Operation().OperationID
	}
	return ""
}