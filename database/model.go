package database

import (
	"gorm.io/gorm"
	"time"
)

type Migrate struct {
	Model any
	Seed  func(db *gorm.DB)
}

type Fields struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
