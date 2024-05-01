package helper

import (
	"fmt"
	"github.com/duxweb/go-fast/database"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"math"
)

type Pagination struct {
	PageSize int   `json:"pageSize,omitempty"`
	Page     int   `json:"page,omitempty"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func NewPagination(page, pageSize int) *Pagination {
	return &Pagination{
		PageSize: pageSize,
		Page:     page,
	}
}

func Paginate(pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64
		db.Session(&gorm.Session{}).Count(&totalRows)
		pagination.Total = totalRows
		totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.PageSize)))
		pagination.Pages = totalPages
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func MorphTo[T any, R any](data []T, model any) {

	hasIds := lo.Uniq(lo.Map[T, uint](data, func(item T, index int) uint {
		return item[ID]
	}))

	userList := []model.SystemUser{}
	err := database.Gorm().Model(model.SystemUser{}).Where("id IN ?", hasIds).Find(&userList).Error
	if err != nil {
		fmt.Println("err", err)
	}
	userData := lo.KeyBy[uint, model.SystemUser](userList, func(item model.SystemUser) uint {
		return item.ID
	})

	data = lo.Map[model.LogOperate, model.LogOperate](data, func(item model.LogOperate, index int) model.LogOperate {
		item.User = userData[item.UserID]
		return item
	})
}

func FormatData[T any](data []T, call func(item T, index int) map[string]any, page *Pagination) ([]map[string]any, map[string]any) {
	meta := map[string]any{}
	transform := lo.Map[T, map[string]any](data, call)
	if page != nil {
		meta["page"] = page.GetOffset()
		meta["total"] = page.Total
		meta["pages"] = page.Pages
	}
	return transform, meta
}
