package models

import (
	"math"

	"gorm.io/gorm"
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
