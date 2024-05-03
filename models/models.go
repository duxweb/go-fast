package models

import (
	"github.com/golang-module/carbon/v2"
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
