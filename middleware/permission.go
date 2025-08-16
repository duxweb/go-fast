package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	duxAuth "github.com/duxweb/go-fast/v2/auth"
	"github.com/duxweb/go-fast/v2/route"
)

type PermissionFun func(id string) ([]string, []string, error)

func PermissionMiddleware(permission PermissionFun) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// 从上下文中获取认证信息
		auth, ok := ctx.Context().Value("auth").(*duxAuth.JwtClaims)
		if !ok {
			huma.WriteErr(nil, ctx, 500, "Permissions must be authorized by the user after")
			return
		}

		// 获取路由名称
		routeName := route.GetRouteName(ctx)
		if routeName == "" {
			next(ctx)
			return
		}

		// 获取权限信息
		allPermission, userPermission, err := permission(auth.ID)
		if err != nil {
			huma.WriteErr(nil, ctx, 500, err.Error())
			return
		}

		// 将权限信息存入上下文
		ctx = huma.WithValue(ctx, "permissions", allPermission)
		ctx = huma.WithValue(ctx, "userPermissions", userPermission)

		// 检查权限
		err = Can(ctx, routeName)
		if err != nil {
			ctx.SetStatus(403)
			return
		}

		next(ctx)
	}
}
