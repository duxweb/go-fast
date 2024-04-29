package action

import (
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/response"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
	"strings"
)

func (t *Resources[T]) List(ctx echo.Context) error {
	var err error
	if t.initFun != nil {
		err = t.initFun(ctx)
		if err != nil {
			return err
		}
	}

	params := map[string]any{}
	err = (&echo.DefaultBinder{}).BindQueryParams(ctx, &params)
	if err != nil {
		return err
	}

	if params["pageSize"] == nil {
		t.Pagination.Status = false
	}

	pageSize := 0
	if t.Pagination.Status {
		pageSize = lo.Ternary[int](params["pageSize"] != nil, cast.ToInt(params["pageSize"]), t.Pagination.PageSize)
	}

	query := database.Gorm().Model(t.Model)

	if params["id"] != nil {
		query = query.Where(t.Key+" = ?", params["id"])
	}

	if t.queryManyFun != nil {
		query = t.queryManyFun(query, params, ctx)
	}

	if t.queryFun != nil {
		query = t.queryFun(query, params, ctx)
	}

	if params["ids"] != nil {
		ids := strings.Split(params["ids"].(string), ",")
		ids = lo.Filter[string](ids, func(item string, index int) bool {
			if item != "" {
				return true
			}
			return false
		})
		query.Where(t.Key+" in ?", ids)
		query.Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(" + t.Key + ",?)", Vars: []any{ids}, WithoutParentheses: true},
		})
	}

	sorts := t.getSorts(params)
	for k, v := range sorts {
		query.Order(k + " " + cast.ToString(v))
	}

	models := make([]T, 0)
	var pagination *helper.Pagination
	if t.Pagination.Status {
		pagination = helper.NewPagination(cast.ToInt(params["page"]), pageSize)
		err = query.Scopes(helper.Paginate(pagination)).Find(&models).Error
	} else {
		err = query.Find(&models).Error
	}

	if err != nil {
		return err
	}

	data := make([]map[string]any, 0)
	meta := map[string]any{}
	if t.transformFun != nil {
		data, meta = helper.FormatData[T](models, t.transformFun, pagination)
	}

	data = t.filterData(data, t.IncludesMany, t.ExcludesMany)

	return response.Send(ctx, response.Data{
		Data: data,
		Meta: meta,
	})
}
