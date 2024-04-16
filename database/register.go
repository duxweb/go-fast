package database

import (
	"github.com/duxweb/go-fast/annotation"
)

func Register(files []*annotation.File) {

	for _, file := range files {
		for _, item := range file.Annotations {
			if item.Name != "AutoMigrate" {
				continue
			}
			GormMigrate(item.Func)
		}
	}

}
