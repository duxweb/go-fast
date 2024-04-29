package action

import (
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type Pagination struct {
	Status   bool
	PageSize int
}

type Resources[T any] struct {
	Model        T
	Key          string
	Tree         bool
	Pagination   Pagination
	IncludesMany []string
	ExcludesMany []string
	IncludesOne  []string
	ExcludesOne  []string
	initFun      InitFun
	transformFun TransformFun[T]
	queryFun     QueryRequestFun
	queryManyFun QueryRequestFun
	queryOneFun  QueryRequestFun
	metaManyFun  MetaManyFun[T]
	metaOneFun   MetaOneFun[T]
	validatorFun ValidatorFun
	formatFun    FormatFun[T]
	ActionList   bool
	ActionShow   bool
	ActionCreate bool
	Extend       map[string]any
}

func New[T any](model T) *Resources[T] {
	return &Resources[T]{
		Key:   "id",
		Tree:  false,
		Model: model,
		Pagination: Pagination{
			Status:   true,
			PageSize: 10,
		},
		IncludesMany: []string{},
		ExcludesMany: []string{},
		IncludesOne:  []string{},
		ExcludesOne:  []string{},
		ActionList:   true,
		ActionShow:   true,
		ActionCreate: true,
		Extend:       map[string]any{},
	}
}

type InitFun func(e echo.Context) error

// Init 初始化回调
func (t *Resources[T]) Init(call InitFun) {
	t.initFun = call
}

type TransformFun[T any] func(item T, index int) map[string]any

// Transform 字段转换
func (t *Resources[T]) Transform(call TransformFun[T]) {
	t.transformFun = call
}

// Query 通用查询
func (t *Resources[T]) Query(call QueryRequestFun) {
	t.queryFun = call
}

type QueryRequestFun func(tx *gorm.DB, params map[string]any, e echo.Context) *gorm.DB

// QueryMany 多条数据查询
func (t *Resources[T]) QueryMany(call QueryRequestFun) {
	t.queryManyFun = call
}

// QueryOne 单条数据查询
func (t *Resources[T]) QueryOne(call QueryRequestFun) {
	t.queryOneFun = call
}

type MetaManyFun[T any] func(orm *gorm.DB, data []T, e echo.Context)

// MetaMany 多条元数据
func (t *Resources[T]) MetaMany(call MetaManyFun[T]) {
	t.metaManyFun = call
}

type MetaOneFun[T any] func(data T, e echo.Context)

// MetaOne 单条元数据
func (t *Resources[T]) MetaOne(call MetaManyFun[T]) {
	t.metaManyFun = call
}

type ValidatorFun func(data map[string]any, e echo.Context) (validator.ValidatorRule, error)

// Validator 数据验证
func (t *Resources[T]) Validator(call ValidatorFun) {
	t.validatorFun = call
}

type FormatFun[T any] func(item T, index int) map[string]any

// Format 数据格式化
func (t *Resources[T]) Format(call FormatFun[T]) {
	t.formatFun = call
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

type Result map[string]func(ctx echo.Context) error

func (t *Resources[T]) Result() Result {
	result := Result{}
	if t.ActionList {
		result["list"] = t.List
	}
	if t.ActionShow {
		result["show"] = t.Show
	}
	if t.ActionCreate {
		result["create"] = t.Create
	}
	return result
}
