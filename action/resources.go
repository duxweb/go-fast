package admin

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type Pagination struct {
	Status   bool
	PageSize uint
}

type Resources[TModel any, TParams any] struct {
	Model        any
	Params       any
	Key          string
	Tree         bool
	Pagination   Pagination
	IncludesMany []string
	ExcludesMany []string
	IncludesOne  []string
	ExcludesOne  []string
	initFun      *InitFun
	transformFun *TransformFun[TModel]
	queryFun     *QueryFun
	queryManyFun *QueryRequestFun
	queryOneFun  *QueryRequestFun
	metaManyFun  *MetaManyFun[TModel]
	metaOneFun   *MetaOneFun[TModel]
	validatorFun *ValidatorFun[TParams]
	formatFun    *FormatFun[TParams]
	ActionList   bool
	ActionShow   bool
	Extend       map[string]any
}

func New[TModel any, TParams any](model any, params any) *Resources[TModel, TParams] {
	return &Resources[TModel, TParams]{
		Key:    "id",
		Tree:   false,
		Model:  model,
		Params: params,
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
		Extend:       map[string]any{},
	}
}

type InitFun func(e echo.Context) error

// Init 初始化回调
func (t *Resources[TModel, TParams]) Init(call *InitFun) {
	t.initFun = call
}

type TransformFun[T any] func(item T) map[string]any

// Transform 字段转换
func (t *Resources[TModel, TParams]) Transform(call *TransformFun[TModel]) {
	t.transformFun = call
}

type QueryFun func(orm *gorm.DB)

// Query 通用查询
func (t *Resources[TModel, TParams]) Query(call *QueryFun) {
	t.queryFun = call
}

type QueryRequestFun func(orm *gorm.DB, e echo.Context)

// QueryMany 多条数据查询
func (t *Resources[TModel, TParams]) QueryMany(call *QueryRequestFun) {
	t.queryManyFun = call
}

// QueryOne 单条数据查询
func (t *Resources[TModel, TParams]) QueryOne(call *QueryRequestFun) {
	t.queryOneFun = call
}

type MetaManyFun[T any] func(orm *gorm.DB, data []T, e echo.Context)

// MetaMany 多条元数据
func (t *Resources[TModel, TParams]) MetaMany(call *MetaManyFun[TModel]) {
	t.metaManyFun = call
}

type MetaOneFun[T any] func(data T, e echo.Context)

// MetaOne 单条元数据
func (t *Resources[TModel, TParams]) MetaOne(call *MetaManyFun[TModel]) {
	t.metaManyFun = call
}

type ValidatorFun[T any] func(data T, e echo.Context)

// Validator 数据验证
func (t *Resources[TModel, TParams]) Validator(call *ValidatorFun[TParams]) {
	t.validatorFun = call
}

type FormatFun[T any] func(params T, e echo.Context) any

// Format 数据格式化
func (t *Resources[TModel, TParams]) Format(call *FormatFun[TParams]) {
	t.formatFun = call
}

// getSorts 获取排序规则
func (t *Resources[TModel, TParams]) getSorts(params map[string]any) map[string]string {
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
