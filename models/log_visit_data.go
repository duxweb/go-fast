package models

import "github.com/dromara/carbon/v2"

// LogVisitData @AutoMigrate()
type LogVisitData struct {
	Fields
	Date     carbon.Date `gorm:"type:date;comment:日期" json:"date"`
	HasType  string      `gorm:"size:250;comment:关联类型" json:"has_type"`
	HasId    uint        `gorm:"size:20;comment:关联 id" json:"has_id"`
	UUID     string      `gorm:"size:250;comment:唯一标识" json:"uuid"`
	Driver   string      `gorm:"size:250;comment:设备" json:"driver"`
	Ip       string      `gorm:"size:100;comment:IP" json:"ip"`
	Browser  string      `gorm:"size:250;comment:浏览器" json:"browser"`
	Country  string      `gorm:"size:100;comment:国家" json:"country"`
	Province string      `gorm:"size:100;comment:省份" json:"province"`
	City     string      `gorm:"size:100;comment:城市" json:"city"`
	Num      uint        `gorm:"size:100;default:0;comment:PV" json:"num"`
}
