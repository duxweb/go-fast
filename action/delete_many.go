package action

import (
	"context"

	"github.com/duxweb/go-fast/v2/resp"
)

// DeleteMany 批量删除记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) DeleteMany(ctx context.Context, input *DeleteManyInput) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
	for _, id := range input.Body {
		if id == "" {
			continue
		}
		deleteInput := &DeleteInput{ID: id}
		_, err := res.Delete(ctx, deleteInput)
		if err != nil {
			return nil, err
		}
	}

	return resp.Send(ctx, resp.Data[any, resp.EmptyMeta]{
		Message: "批量删除成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}
