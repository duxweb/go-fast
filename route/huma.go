package route

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/views"
)

// CommonError 通用错误结构，实现 huma.StatusError 接口
// CommonError implements huma.StatusError interface for unified response format
type CommonError struct {
	status  int
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta"`
}

// Error 返回错误消息
// Error returns the error message
func (e *CommonError) Error() string {
	return e.Message
}

// GetStatus 返回 HTTP 状态码
// GetStatus returns the HTTP status code
func (e *CommonError) GetStatus() int {
	return e.status
}

// NewHuma 创建 Huma API 实例，配置统一的响应格式和安全方案
// NewHuma creates a Huma API instance with unified response format and security schemes
func NewHuma(name string, prefix string) huma.API {

	// 覆盖默认的 huma.NewError 函数，使用我们的自定义错误格式
	// Override default huma.NewError function to use our custom error format
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		var data any = nil
		if len(errs) > 0 {
			errorDetails := make([]map[string]any, len(errs))
			for i, err := range errs {
				errorDetails[i] = map[string]any{
					"message": err.Error(),
				}
			}
			data = errorDetails
		}

		return &CommonError{
			status:  status,
			Code:    status,
			Message: message,
			Data:    data,
			Meta:    map[string]any{},
		}
	}

	config := huma.DefaultConfig(name+" API", global.Version)

	config.OpenAPIPath = prefix + "/openapi"
	config.SchemasPath = prefix + "/schemas"
	config.DocsPath = prefix + "/docs"

	// 配置安全方案
	// Configure security schemes
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT 认证，格式：Bearer {token}",
		},
		"apiSignature": {
			Type:        "apiKey",
			In:          "header",
			Name:        "Content-MD5",
			Description: "API 签名认证。需要在请求头中包含：AccessKey（访问密钥ID）、Content-Date（Unix时间戳）、Content-MD5（签名）。签名算法：SHA256_HMAC(path + '\\n' + queryString + '\\n' + timestamp, appSecret)。时间戳有效期：10秒",
			Extensions: map[string]any{
				"x-signature-headers": []string{
					"AccessKey",
					"Content-Date",
					"Content-MD5",
				},
				"x-signature-algorithm": "SHA256_HMAC",
				"x-signature-example": map[string]any{
					"accessKey":   "app_123456",
					"appSecret":   "your_app_secret",
					"timestamp":   "1704067200",
					"path":        "/api/v1/users",
					"queryString": "page=1&limit=10",
					"signData":    "path + '\\n' + queryString + '\\n' + timestamp",
					"signature":   "SHA256_HMAC(signData, appSecret) -> hex",
				},
				"x-signature-timeout": "10 seconds",
			},
		},
	}

	// 设置全局安全要求 - 所有接口默认需要 bearer token
	// Set global security requirements - all endpoints require bearer token by default
	config.Security = []map[string][]string{
		{"bearer": {}},
	}

	// 配置服务器
	// Configure servers
	config.Servers = []*huma.Server{
		{
			URL:         fmt.Sprintf("http://%s:%s", global.Web.Host, global.Web.Port),
			Description: "Local server",
		},
		{
			URL:         global.Web.Domain,
			Description: "Domain server",
		},
	}

	// 定义统一的错误响应 Schema
	// Define unified error response schema
	errorSchema := &huma.Schema{
		Type: "object",
		Properties: map[string]*huma.Schema{
			"code": {
				Type:     "integer",
				Title:    "状态码",
				Default:  500,
				Examples: []any{500},
			},
			"message": {
				Type:     "string",
				Title:    "错误消息",
				Examples: []any{"Bad Request"},
			},
			"data": {
				Title: "错误数据",
			},
			"meta": {
				Type:    "object",
				Title:   "元数据",
				Default: map[string]any{},
			},
		},
		Required: []string{"code", "message", "data", "meta"},
	}

	// 设置统一的错误响应格式
	// Set unified error response format
	config.Components.Responses = map[string]*huma.Response{
		"ErrorResponse": {
			Description: "统一错误响应格式",
			Content: map[string]*huma.MediaType{
				"application/json": {
					Schema: errorSchema,
				},
			},
		},
	}

	return humaecho.New(global.Router, config)
}

// Get 注册一个 GET 操作
// Get registers a GET operation.
func Get[I any, O any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	options.Method = "GET"
	return Add(s, http.MethodGet, path, name, options, handler)
}

// Post 注册一个 POST 操作
// Post registers a POST operation.
func Post[I any, O any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	options.Method = "POST"
	return Add(s, http.MethodPost, path, name, options, handler)
}

// Put 注册一个 PUT 操作
// Put registers a PUT operation.
func Put[I any, O any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	options.Method = "PUT"
	return Add(s, http.MethodPut, path, name, options, handler)
}

// Delete 注册一个 DELETE 操作
// Delete registers a DELETE operation.
func Delete[I any, O any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	options.Method = "DELETE"
	return Add(s, http.MethodDelete, path, name, options, handler)
}

// Head 注册一个 HEAD 操作
// Head registers an HEAD operation.
func Head[I any, O any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	options.Method = "HEAD"
	return Add(s, http.MethodOptions, path, name, options, handler)
}

// Patch 注册一个 PATCH 操作
// Patch registers a PATCH operation.
func Patch[I any, O any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	options.Method = "PATCH"
	return Add(s, http.MethodPatch, path, name, options, handler)
}

// Add 使用 Huma 注册一个路由操作
// Add registers an operation with Huma.
// 对可携带请求体的方法，输入会被 Body 包装
// For methods that can have a body, the input is wrapped into a Body struct.
func Add[I any, O any](s *RouterData, method string, path string, name string, options huma.Operation, handler func(context.Context, *I) (*O, error)) *RouterItem {
	fullName := s.Name + "." + name
	options.OperationID = fullName
	options.Path = path
	options.Method = method

	// 添加默认的错误响应（如果没有定义）
	// Add default error responses if not defined
	if options.Responses == nil {
		options.Responses = map[string]*huma.Response{}
	}

	// 定义错误响应 schema
	errorSchema := &huma.Schema{
		Type: "object",
		Properties: map[string]*huma.Schema{
			"code": {
				Type:  "integer",
				Title: "状态码",
			},
			"message": {
				Type:  "string",
				Title: "错误消息",
			},
			"data": {
				Title: "错误数据",
			},
			"meta": {
				Type:  "object",
				Title: "元数据",
			},
		},
		Required: []string{"code", "message", "data", "meta"},
	}

	// 添加常见的错误状态码响应
	errorStatuses := map[string]string{
		"400": "请求参数错误",
		"401": "未授权",
		"403": "禁止访问",
		"404": "资源不存在",
		"422": "参数验证失败",
		"500": "服务器内部错误",
	}

	for statusCode, description := range errorStatuses {
		if _, exists := options.Responses[statusCode]; !exists {
			options.Responses[statusCode] = &huma.Response{
				Description: description,
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: errorSchema,
					},
				},
			}
		}
	}

	// 使用路由组自己的 Huma 实例而不是全局实例
	// Use route group's own Huma instance instead of global one
	if s.HumaRouter != nil {
		huma.Register(s.HumaRouter, options, handler)
	}

	item := RouterItem{
		Method: method,
		Path:   path,
		Name:   fullName,
	}

	s.Data = append(s.Data, &item)

	return &item
}

// Page 注册一个返回 HTML 页面的操作
// Page registers an operation that returns HTML page.
func Page[I any](s *RouterData, path string, name string, options huma.Operation, handler func(context.Context, *I) (*views.HTMLOutput, error)) *RouterItem {
	options.Method = "GET"
	return Add[I, views.HTMLOutput](s, http.MethodGet, path, name, options, handler)
}

// UseMiddleware 为路由组添加 Huma 中间件
// UseMiddleware adds Huma middleware to the router group.
func (s *RouterData) UseMiddleware(middlewares ...func(huma.Context, func(huma.Context))) {
	for _, middleware := range middlewares {
		s.HumaRouter.UseMiddleware(middleware)
	}
}
