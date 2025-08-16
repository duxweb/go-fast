package action

import (
	"context"

	"github.com/duxweb/go-fast/v2/resp"
)

// TrashMany 批量彻底删除记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) TrashMany(ctx context.Context, input *TrashManyInput) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
	for _, id := range input.IDs {
		if id == "" {
			continue
		}
		trashInput := &DeleteInput{ID: id}
		_, err := res.Trash(ctx, trashInput)
		if err != nil {
			return nil, err
		}
	}

	return resp.Send(ctx, resp.Data[any, resp.EmptyMeta]{
		Message: "批量彻底删除成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}