package resp

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/i18n"
	"github.com/duxweb/go-fast/v2/views"
)

// EmptyMeta 空元数据
type EmptyMeta struct{}

// Render 使用 views 系统渲染模板并返回 HTML 响应
func Render(ctx context.Context, templateName string, data any) (*views.HTMLOutput, error) {
	var buf bytes.Buffer

	// 使用现有的 View.Render 方法
	err := views.View.Render(&buf, templateName, data, ctx)
	if err != nil {
		return nil, err
	}

	return views.NewHTMLResponse(buf.String()), nil
}

type Data[T, M any] struct {
	Code        int    `json:"code" example:"200"`
	Message     string `json:"message" example:"ok"`
	MessageLang string `json:"-"`
	Data        T      `json:"data"`
	Meta        M      `json:"meta"`
}

func RawSend(ctx huma.Context, data Data[any, any], code ...int) error {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}
	if data.Message == "" {
		data.Message = "ok"
	}
	if data.MessageLang != "" {
		data.Message = i18n.T(ctx.Context(), data.MessageLang)
	}
	if data.Meta == nil {
		data.Meta = map[string]any{}
	}
	if data.Code == 0 {
		data.Code = statusCode
	}

	// 直接使用 huma.Context 输出
	ctx.SetStatus(statusCode)
	ctx.SetHeader("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ctx.BodyWriter().Write(jsonData)
	return err
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
