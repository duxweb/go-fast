package action

import (
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	coreModel "github.com/duxweb/go-fast/models"
	"github.com/duxweb/go-fast/response"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gorm.io/gorm/clause"
	"strings"
)

func (t *Resources[T]) List(ctx *fiber.Ctx) error {
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

	if t.Tree || !params.Get("pageSize").Exists() {
		t.Pagination.Status = false
	}

	pageSize := 0
	if t.Pagination.Status {
		pageSize = lo.Ternary[int](params.Get("pageSize").Exists(), int(params.Get("pageSize").Uint()), t.Pagination.PageSize)
	}

	query := database.Gorm().Model(t.Model).Debug()

	if params.Get("id").Exists() {
		query = query.Where(t.Key+" = ?", params.Get("id"))
	}

	if params.Get("ids").Exists() {
		ids := strings.Split(params.Get("ids").String(), ",")
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

	if t.queryManyFun != nil {
		query = t.queryManyFun(query, params, ctx)
	}

	if t.queryFun != nil {
		query = t.queryFun(query, ctx)
	}

	models := make([]T, 0)
	var pagination *coreModel.Pagination
	if t.Pagination.Status {
		pagination = coreModel.NewPagination(int(params.Get("page").Int()), pageSize)
		err = query.Scopes(coreModel.Paginate(pagination)).Find(&models).Error
	} else {
		if t.Tree {
			query = query.Preload(clause.Associations, coreModel.ChildrenPreload).Where("parent_id = 0")
		}
		err = query.Find(&models).Error
	}
	if err != nil {
		return err
	}

	if t.manyAfterFun != nil {
		models = t.manyAfterFun(models, params, ctx)
	}

	data := make([]map[string]any, 0)

	meta := map[string]any{}
	if t.TransformFun != nil {
		data, meta = coreModel.FormatData[T](models, t.TransformFun, pagination)
	}

	data = t.filterData(data, t.IncludesMany, t.ExcludesMany)

	return response.Send(ctx, response.Data{
		Data: data,
		Meta: meta,
	})
}

type ManyCallFun[T any] func(data []T, params *gjson.Result, ctx *fiber.Ctx) []T

func (t *Resources[T]) ManyAfter(call ManyCallFun[T]) {
	t.manyAfterFun = call
}

// getSorts 获取排序规则
func (t *Resources[T]) getSorts(params *gjson.Result) map[string]string {
	data := map[string]string{}
	for key, value := range params.Map() {
		if !strings.HasSuffix(key, "_sort") {
			continue
		}
		if value.String() != "asc" && value.String() != "desc" {
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
