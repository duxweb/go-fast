package resp

import (
	"context"

	"github.com/duxweb/go-fast/v2/i18n"
	"github.com/labstack/echo/v4"
)

// EmptyMeta 空元数据
type EmptyMeta struct{}

func RawRender(ctx echo.Context, app string, name string, bind any, code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	templateName := name
	if app != "" {
		templateName = app + ":" + name
	}

	return ctx.Render(statusCode, templateName, bind)
}

type Data[T, M any] struct {
	Code        int    `json:"code" example:"200"`
	Message     string `json:"message" example:"ok"`
	MessageLang string `json:"-"`
	Data        T      `json:"data"`
	Meta        M      `json:"meta"`
}

func RawSend(ctx echo.Context, data Data[any, any], code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	if data.Message == "" {
		data.Message = "ok"
	}
	if data.MessageLang != "" {
		data.Message = i18n.T(ctx.Request().Context(), data.MessageLang)
	}
	if data.Meta == nil {
		data.Meta = echo.Map{}
	}
	if data.Code == 0 {
		data.Code = statusCode
	}
	return ctx.JSON(statusCode, data)
}

// HumaResponse Huma API 统一响应格式
// HumaResponse unified response format for Huma API
type HumaResponse[T, M any] struct {
	StatusCode int `status:"200"`
	Body       struct {
		Code    int    `json:"code" example:"200" doc:"状态码"`
		Message string `json:"message" example:"ok" doc:"消息"`
		Data    T      `json:"data" doc:"数据"`
		Meta    M      `json:"meta" doc:"元数据"`
	}
}

// NewHumaResponse 统一响应处理
// NewHumaResponse send unified response
func Send[T, M any](ctx context.Context, data Data[T, M], code ...int) *HumaResponse[T, M] {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	if data.Message == "" {
		data.Message = "ok"
	}
	if data.MessageLang != "" {
		data.Message = i18n.T(ctx, data.MessageLang)
	}
	if data.Code == 0 {
		data.Code = statusCode
	}

	result := &HumaResponse[T, M]{}
	result.StatusCode = statusCode

	result.Body.Code = data.Code
	result.Body.Message = data.Message
	result.Body.Data = data.Data
	result.Body.Meta = data.Meta

	return result
}
