package database

import (
	"github.com/duxweb/go-fast/v2/annotation"
	"github.com/duxweb/go-fast/v2/models"
)

func Register() {
	GormMigrate(models.LogOperate{}, models.LogLogin{})
	for _, file := range annotation.Annotations {
		for _, item := range file.Annotations {
			if item.Name != "AutoMigrate" {
				continue
			}
			if item.Func == nil {
				panic("database func not set: " + file.Name)
			}
			GormMigrate(item.Func)
		}
	}

}
