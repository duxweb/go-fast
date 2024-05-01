package action

import (
	"errors"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (t *Resources[T]) Show(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(ctx)
		if err != nil {
			return err
		}
	}

	params := map[string]any{}
	err = ctx.Bind(&params)
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
	if t.transformFun != nil && !isEmpty {
		data = t.transformFun(model, 0)
	}

	filterData := t.filterData([]map[string]any{data}, t.IncludesMany, t.ExcludesMany)
	if len(filterData) > 0 {
		data = filterData[0]
	}

	return response.Send(ctx, response.Data{
		Data: data,
		Meta: meta,
	})
}

type OneCallFun[T any] func(data T, params map[string]any, ctx echo.Context) T

func (t *Resources[T]) OneAfter(call OneCallFun[T]) {
	t.oneAfterFun = call
}

func (t *Resources[T]) getOne(ctx echo.Context, model *T, id string, params map[string]any) error {
	query := database.Gorm().Unscoped().Model(t.Model).Where(t.Key+" = ?", id)
	if t.queryOneFun != nil {
		query = t.queryOneFun(query, params, ctx)
	}
	if t.queryFun != nil {
		query = t.queryFun(query, params, ctx)
	}
	return query.First(model).Error
}
