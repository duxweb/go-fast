package web

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/duxweb/go-fast/v2/resp"
	"github.com/duxweb/go-fast/v2/route"
)

// TestInput 测试输入结构
type TestInput struct {
	ID   string `path:"id" example:"1" doc:"测试ID"`
	Name string `json:"name" example:"test" doc:"测试名称"`
}

// TestData 测试数据结构
type TestData struct {
	ID   string `json:"id" example:"1" doc:"测试ID"`
	Name string `json:"name" example:"test" doc:"测试名称"`
}

// Test @Route()
func Test() {
	group := route.GetRouter("web")

	route.Get(group, "/test/{id}", "test", huma.Operation{
		OperationID: "getTest",
		Summary:     "获取测试数据",
		Description: "根据ID获取测试数据",
		Tags:        []string{"test"},
	}, func(ctx context.Context, input *TestInput) (*resp.HumaResponse[TestData, resp.EmptyMeta], error) {
		data := TestData{
			ID:   input.ID,
			Name: input.Name,
		}

		return resp.Send(ctx, resp.Data[TestData, resp.EmptyMeta]{
			Data: data,
			Meta: resp.EmptyMeta{},
		}), nil
	})
}
