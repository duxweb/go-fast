package action

import (
	"context"

	"github.com/duxweb/go-fast/v2/resp"
)

// RestoreMany 批量恢复软删除记录方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) RestoreMany(ctx context.Context, input *RestoreManyInput) (*resp.HumaResponse[any, resp.EmptyMeta], error) {
	for _, id := range input.IDs {
		if id == "" {
			continue
		}
		restoreInput := &DeleteInput{ID: id}
		_, err := res.Restore(ctx, restoreInput)
		if err != nil {
			return nil, err
		}
	}

	return resp.Send(ctx, resp.Data[any, resp.EmptyMeta]{
		Message: "批量恢复成功",
		Data:    nil,
		Meta:    resp.EmptyMeta{},
	}), nil
}