package action

import (
	"github.com/duxweb/go-fast/validator"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type Pagination struct {
	Status   bool
	PageSize int
}

type Resources[T any] struct {
	Model            T
	Key              string
	Tree             bool
	Pagination       Pagination
	IncludesMany     []string
	ExcludesMany     []string
	IncludesOne      []string
	ExcludesOne      []string
	queryParams      any
	initFun          InitFun
	transformFun     TransformFun[T]
	queryFun         QueryFun
	queryManyFun     QueryRequestFun
	queryOneFun      QueryRequestFun
	metaManyFun      MetaManyFun[T]
	metaOneFun       MetaOneFun[T]
	manyAfterFun     ManyCallFun[T]
	oneAfterFun      OneCallFun[T]
	validatorFun     ValidatorFun
	formatFun        FormatFun[T]
	createBeforeFun  ActionCallParamsFun[T]
	createAfterFun   ActionCallParamsFun[T]
	editBeforeFun    ActionCallParamsFun[T]
	editAfterFun     ActionCallParamsFun[T]
	saveBeforeFun    ActionCallParamsFun[T]
	saveAfterFun     ActionCallParamsFun[T]
	storeBeforeFun   ActionCallParamsFun[T]
	storeAfterFun    ActionCallParamsFun[T]
	deleteBeforeFun  ActionCallFun[T]
	deleteAfterFun   ActionCallFun[T]
	trashBeforeFun   ActionCallFun[T]
	trashAfterFun    ActionCallFun[T]
	restoreBeforeFun ActionCallFun[T]
	restoreAfterFun  ActionCallFun[T]
	ActionList       bool
	ActionShow       bool
	ActionCreate     bool
	ActionEdit       bool
	ActionDelete     bool
	ActionStore      bool
	ActionSoftDelete bool
	Extend           map[string]any
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
		IncludesMany:     []string{},
		ExcludesMany:     []string{},
		IncludesOne:      []string{},
		ExcludesOne:      []string{},
		ActionList:       true,
		ActionShow:       true,
		ActionCreate:     true,
		ActionEdit:       true,
		ActionDelete:     true,
		ActionStore:      true,
		ActionSoftDelete: false,
		Extend:           map[string]any{},
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

func (t *Resources[T]) QueryParams(data any) {
	t.queryParams = data
}

type QueryFun func(tx *gorm.DB, e echo.Context) *gorm.DB

// Query 通用查询
func (t *Resources[T]) Query(call QueryFun) {
	t.queryFun = call
}

type QueryRequestFun func(tx *gorm.DB, params *gjson.Result, e echo.Context) *gorm.DB

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

type ValidatorFun func(data *gjson.Result, e echo.Context) (validator.ValidatorRule, error)

// Validator 数据验证
// Docs github.com/go-playground/validator/v10
func (t *Resources[T]) Validator(call ValidatorFun) {
	t.validatorFun = call
}

type FormatFun[T any] func(model *T, data *gjson.Result, e echo.Context) error

// Format 数据格式化
func (t *Resources[T]) Format(call FormatFun[T]) {
	t.formatFun = call
}

type ActionCallParamsFun[T any] func(data *T, params *gjson.Result) error

type ActionCallFun[T any] func(data *T) error

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
	if t.ActionEdit {
		result["edit"] = t.Edit
	}
	if t.ActionDelete {
		result["delete"] = t.Delete
		result["deleteMany"] = t.DeleteMany
	}
	if t.ActionStore {
		result["store"] = t.Store
	}
	if t.ActionSoftDelete {
		result["trash"] = t.Trash
		result["trashMany"] = t.TrashMany
		result["restore"] = t.Restore
		result["restoreMany"] = t.RestoreMany
	}
	return result
}
