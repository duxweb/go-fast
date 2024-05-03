package models

import (
	"github.com/golang-module/carbon/v2"
)

// LogVisitSpider @AutoMigrate()
type LogVisitSpider struct {
	Fields
	Date    carbon.Date `gorm:"type:date;comment:日期" json:"date"`
	HasType string      `gorm:"size:250;comment:关联类型" json:"has_type"`
	HasId   uint        `gorm:"size:20;comment:关联 id" json:"has_id"`
	Name    string      `gorm:"size:250;comment:蜘蛛名" json:"name"`
	Path    string      `gorm:"size:250;comment:页面路径" json:"path"`
	Num     uint        `gorm:"default:0;comment:访客量" json:"num"`
}
