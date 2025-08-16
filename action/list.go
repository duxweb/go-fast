package action

import (
	"context"
	"reflect"
	"strings"

	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/helper"
	coreModel "github.com/duxweb/go-fast/v2/models"
	"github.com/duxweb/go-fast/v2/resp"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// List 列表查询方法
func (res *Resources[Model, Info, Params, Data, ListMeta, DetailMeta]) List(ctx context.Context, input *Params) (*resp.HumaResponse[[]Info, ListMeta], error) {

	params := helper.StructToGJson(input)

	pageStatus := res.Pagination.Status

	if res.Tree || !params.Get("pageSize").Exists() {
		pageStatus = false
	}

	pageSize := 0
	if pageStatus {
		pageSize = lo.Ternary[int](params.Get("pageSize").Exists(), int(params.Get("pageSize").Uint()), res.Pagination.PageSize)
	}

	// 获取模型类型
	var model Model
	if res.model != nil {
		// 使用反射创建模型实例
		modelType := reflect.TypeOf(res.model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		model = reflect.New(modelType).Interface().(Model)
	}

	query := database.Gorm().Model(model).Debug()

	if params.Get("id").Exists() {
		query = query.Where(res.Key+" = ?", params.Get("id"))
	}

	if params.Get("ids").Exists() {
		ids := strings.Split(params.Get("ids").String(), ",")
		ids = lo.Filter[string](ids, func(item string, index int) bool {
			if item != "" {
				return true
			}
			return false
		})

		query = query.Where(res.Key+" in ?", ids).Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(" + res.Key + ",?)", Vars: []any{ids}, WithoutParentheses: true},
		})
	}

	if res.Tree {
		query = query.Set("tree_sort", res.TreeSort)
		if res.TreeSort != "" {
			query = query.Order(res.TreeSort + " ASC")
		}
	}

	sorts := getSorts(params)
	for k, v := range sorts {
		query = query.Order(k + " " + cast.ToString(v))
	}

	// 应用筛选查询
	if res.filterFun != nil {
		filterFn := res.filterFun.(func(*gorm.DB, *Params) *gorm.DB)
		query = filterFn(query, input)
	}

	// 应用全局查询
	if res.queryFun != nil {
		queryFn := res.queryFun.(func(*gorm.DB, context.Context) *gorm.DB)
		query = queryFn(query, ctx)
	}

	if res.preload != nil {
		for _, v := range res.preload {
			query = query.Preload(v)
		}
	}

	models := make([]Model, 0)
	var pagination *coreModel.Pagination
	if pageStatus {
		pagination = coreModel.NewPagination(int(params.Get("page").Int()), pageSize)
		err := query.Scopes(coreModel.Paginate(pagination)).Find(&models).Error
		if err != nil {
			return nil, err
		}
	} else {
		if res.Tree {
			config := coreModel.TreePreloadConfig{
				Sort:     res.TreeSort,
				Preloads: res.preload,
			}
			query = query.Preload("Children", coreModel.TreePreload(config)).Where("parent_id IS NULL")
		}
		err := query.Find(&models).Error
		if err != nil {
			return nil, err
		}
	}

	data := make([]Info, 0)

	if res.transformFun != nil {
		transformFn := res.transformFun.(func(*Model, int, *TransformContext) Info)
		for i, model := range models {
			transformedData := transformFn(&model, i, &TransformContext{IsList: true})
			data = append(data, transformedData)
		}
	}

	// 处理meta数据
	var metaStruct ListMeta
	
	// 如果有自定义meta函数，使用自定义meta
	if res.metaManyFun != nil {
		if metaFn, ok := res.metaManyFun.(func([]Model) ListMeta); ok {
			metaStruct = metaFn(models)
		}
	} else {
		// 使用默认Meta类型，需要转换为ListMeta
		defaultMeta := Meta{}
		if pagination != nil {
			defaultMeta.Total = int(pagination.Total)
			defaultMeta.Page = pagination.Page
			defaultMeta.Limit = pageSize
		}
		
		// 将默认Meta转换为ListMeta（假设ListMeta包含Meta字段或兼容）
		if m, ok := any(defaultMeta).(ListMeta); ok {
			metaStruct = m
		}
	}

	return resp.Send(ctx, resp.Data[[]Info, ListMeta]{
		Data: data,
		Meta: metaStruct,
	}), nil
}

// getSorts 获取排序规则
func getSorts(params *gjson.Result) map[string]string {
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
