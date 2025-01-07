package action

import (
	"errors"

	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

func (t *Resources[T]) Show(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(t, ctx)
		if err != nil {
			return err
		}
	}

	params, err := helper.Qs(ctx)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	var model T
	err = t.getOne(ctx, &model, id, params)
	isEmpty := false
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isEmpty = true
		} else {
			return err
		}
	}

	if t.oneAfterFun != nil {
		model = t.oneAfterFun(model, params, ctx)
	}

	data := map[string]any{}
	meta := map[string]any{}
	if t.TransformFun != nil && !isEmpty {
		data = t.TransformFun(&model, 0)
	}

	filterData := t.filterData([]map[string]any{data}, t.IncludesMany, t.ExcludesMany)
	if len(filterData) > 0 {
		data = filterData[0]
	}
	if t.metaOneFun != nil {
		mayMeta := t.metaOneFun(model, ctx)
		meta = lo.Assign(meta, mayMeta)
	}

	return response.Send(ctx, response.Data{
		Data: data,
		Meta: meta,
	})
}

type OneCallFun[T any] func(data T, params *gjson.Result, ctx echo.Context) T

func (t *Resources[T]) OneAfter(call OneCallFun[T]) {
	t.oneAfterFun = call
}

func (t *Resources[T]) getOne(ctx echo.Context, model *T, id string, params *gjson.Result) error {
	query := database.Gorm().Unscoped().Model(t.Model).Where(t.Key+" = ?", id)
	if t.queryOneFun != nil {
		query = t.queryOneFun(query, params, ctx)
	}
	if t.queryFun != nil {
		query = t.queryFun(query, ctx)
	}
	if t.preload != nil {
		for _, v := range t.preload {
			query = query.Preload(v)
		}
	}
	return query.First(model).Error
}
