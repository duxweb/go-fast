package admin

import (
	"github.com/duxweb/go-fast/database"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
	"strings"
)

func (t *Resources[TModel, TParams]) List(ctx echo.Context) error {

	var err error
	if t.initFun != nil {
		initFun := *t.initFun
		err = initFun(ctx)
		if err != nil {
			return err
		}
	}

	params := map[string]any{}
	if err = ctx.Bind(&params); err != nil {
		return err
	}

	if params["pageSize"] == nil {
		t.Pagination.Status = false
	}

	var limit uint = 0
	if t.Pagination.Status {
		limit = lo.Ternary[uint](params["pageSize"] != nil, cast.ToUint(params["pageSize"]), t.Pagination.PageSize)
	}

	query := database.Gorm().Model(t.Model)

	if params["id"] != nil {
		query.Where(t.Key+" = ?", params["id"])
	}

	if t.queryManyFun != nil {
		queryManyFun := *t.queryManyFun
		queryManyFun(query, ctx)
	}

	if t.queryFun != nil {
		queryFun := *t.queryFun
		queryFun(query)
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
	return nil
}
