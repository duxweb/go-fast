package models

import "github.com/duxweb/go-fast/v2/models"

// Config @AutoMigrate()
type Book struct {
	models.Fields
	Name    string `gorm:"size:250" json:"name"`
	Content string `gorm:"type:text" json:"value"`
}
