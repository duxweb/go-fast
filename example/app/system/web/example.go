package web

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/resp"
	"github.com/duxweb/go-fast/v2/route"
	"github.com/duxweb/go-fast/v2/validator"
)

// ExampleInput 示例输入 - 展示验证规则
type ExampleInput struct {
	Name     string `json:"name" minLength:"2" maxLength:"50" required:"true" doc:"姓名"`
	Email    string `json:"email" format:"email" required:"true" doc:"邮箱"`
	Phone    string `json:"phone,omitempty" pattern:"^1[3-9]\\d{9}$" doc:"手机号" message:"请输入有效的手机号"`
	Age      int    `json:"age" minimum:"18" maximum:"100" doc:"年龄"`
	Password string `json:"password" minLength:"6" doc:"密码"`
}

// Resolve 实现验证
func (e *ExampleInput) Resolve(ctx huma.Context) []error {
	return validator.ValidateWithContext(ctx.Context(), e)
}

// ExampleOutput 示例输出
type ExampleOutput struct {
	ID      string `json:"id" doc:"ID"`
	Name    string `json:"name" doc:"姓名"`
	Message string `json:"message" doc:"消息"`
}

// Example 注册示例路由
func Example() {
	group := route.GetRouter("web")

	route.Post(group, "/example", "createExample", huma.Operation{
		OperationID: "createExample",
		Summary:     "创建示例",
		Tags:        []string{"example"},
	}, func(ctx context.Context, input *ExampleInput) (*resp.HumaResponse[ExampleOutput, resp.EmptyMeta], error) {
		return resp.Send(ctx, resp.Data[ExampleOutput, resp.EmptyMeta]{
			Data: ExampleOutput{
				ID:      "123",
				Name:    input.Name,
				Message: "创建成功",
			},
			Meta: resp.EmptyMeta{},
		}), nil
	})
}
