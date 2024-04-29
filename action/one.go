package action

import (
	"errors"
	"fmt"
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

	id := ctx.Param("id")

	params := map[string]any{}
	err = (&echo.DefaultBinder{}).BindQueryParams(ctx, &params)
	if err != nil {
		return err
	}

	query := database.Gorm().Debug().Model(t.Model).Where(t.Key+" = ?", id)

	if t.queryOneFun != nil {
		query = t.queryOneFun(query, params, ctx)
	}

	if t.queryFun != nil {
		query = t.queryFun(query, params, ctx)
	}

	isEmpty := false
	var model T
	if err = query.First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isEmpty = true
		} else {
			return err
		}
	}

	fmt.Println("model", model)

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
