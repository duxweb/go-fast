package models

import (
	"github.com/dromara/carbon/v2"
	"gorm.io/gorm"
)

type Fields struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

type Timestamps struct {
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

type SoftDeletes struct {
	DeletedAt carbon.DateTime `json:"deleted_at" gorm:"index"`
}

func FormatData[T any](data []T, call func(item *T, index int) map[string]any, page *Pagination) ([]map[string]any, map[string]any) {
	meta := map[string]any{}
	transform := make([]map[string]any, 0)
	for i, item := range data {
		transform = append(transform, call(&item, i))
	}
	if page != nil {
		meta["page"] = page.GetOffset()
		meta["total"] = page.Total
		meta["pages"] = page.Pages
	}
	return transform, meta
}

func GetTableName(db *gorm.DB, model any) string {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(&model)
	return stmt.Schema.Table
}
