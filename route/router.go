package route

import (
	pathutil "path"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// RouterData 路由分组，统一使用 Huma 管理
// RouterData describes a logical group of routes managed by Huma.
type RouterData struct {
	Name       string
	Prefix     string
	Data       []*RouterItem
	Groups     []*RouterData
	HumaRouter huma.API
}

type RouterItem struct {
	Method string
	Path   string
	Name   string
}

// RouterMiddle 分组中间件类型（Huma）
// RouterMiddle is a Huma middleware for grouping.
type RouterMiddle = func(huma.Context, func(huma.Context))

// New 创建一个带可选中间件链的路由分组，使用 Huma API 实例
// New creates a route group with Huma API instance and optional Huma middleware.
func New(name string, prefix string, middle ...RouterMiddle) *RouterData {
	// 直接使用 huma.New 创建 API 实例
	huma := NewHuma(name, prefix)

	// 应用 Huma 中间件
	for _, middleware := range middle {
		huma.UseMiddleware(middleware)
	}

	return &RouterData{
		Prefix:     prefix,
		Name:       name,
		HumaRouter: huma,
	}
}

// Group 基于父分组创建子分组，使用 Huma 分组功能
// Group creates a subgroup using Huma's grouping functionality.
func Group(s *RouterData, prefix string, name string, middle ...RouterMiddle) *RouterData {
	// 使用 Huma 分组，自动继承父级中间件
	humGrp := huma.NewGroup(s.HumaRouter, prefix)

	// 应用 Huma 中间件
	for _, middleware := range middle {
		humGrp.UseMiddleware(middleware)
	}

	group := &RouterData{
		Prefix:     prefix,
		HumaRouter: humGrp,
		Name:       s.Name + "." + name,
	}
	s.Groups = append(s.Groups, group)
	return group
}

// joinPath 安全拼接分组前缀与本地路径
// joinPath safely concatenates group prefixes and local paths.
func joinPath(prefix, path string) string {
	if prefix == "" {
		if path == "" {
			return "/"
		}
		if strings.HasPrefix(path, "/") {
			return path
		}
		return "/" + path
	}
	// 确保单一前导斜杠，避免出现双分隔符
	// Ensure single leading slash, prevent double separators
	return pathutil.Clean("/" + strings.TrimSuffix(prefix, "/") + "/" + strings.TrimPrefix(path, "/"))
}

func (t *RouterData) ParseTree(prefix string) map[string]any {
	var all []any
	for _, datum := range t.Data {
		all = append(all, map[string]any{
			"name":   datum.Name,
			"method": datum.Method,
			"path":   prefix + datum.Path,
		})
	}
	for _, item := range t.Groups {
		gpath := prefix + item.Prefix
		all = append(all, item.ParseTree(gpath))
	}
	return map[string]any{
		"path": prefix,
		"data": all,
	}
}

func (t *RouterData) ParseData(prefix string) []map[string]any {
	var all []map[string]any
	for _, datum := range t.Data {
		all = append(all, map[string]any{
			"name":   datum.Name,
			"method": datum.Method,
			"path":   prefix + datum.Path,
		})
	}
	for _, item := range t.Groups {
		gpath := prefix + item.Prefix
		data := item.ParseData(gpath)
		all = append(all, data...)
	}
	return all
}
