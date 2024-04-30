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

		query = query.Where(t.Key+" in ?", ids).Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(" + t.Key + ",?)", Vars: []any{ids}, WithoutParentheses: true},
		})
	}

	sorts := t.getSorts(params)
	for k, v := range sorts {
		query = query.Order(k + " " + cast.ToString(v))
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

// getSorts 获取排序规则
func (t *Resources[T]) getSorts(params map[string]any) map[string]string {
	data := map[string]string{}
	for key, value := range params {
		if !strings.HasSuffix(key, "_sort") {
			continue
		}
		if value != "asc" && value != "desc" {
			continue
		}
		field := key[0 : len(key)-5]
		data[field] = cast.ToString(value)
	}
	return data
}

func (t *Resources[T]) filterData(data []map[string]any, includes []string, excludes []string) []map[string]any {
	result := make([]map[string]any, 0)
	for _, item := range data {
		datum := item
		if len(includes) > 0 {
			datum = lo.PickBy[string, any](item, func(key string, value any) bool {
				return lo.IndexOf[string](includes, key) != -1
			})
		}
		if len(excludes) > 0 {
			datum = lo.PickBy[string, any](datum, func(key string, value any) bool {
				return lo.IndexOf[string](excludes, key) == -1
			})
		}
		result = append(result, datum)
	}
	return result
}
