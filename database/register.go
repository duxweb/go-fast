package database

import "github.com/duxweb/go-fast/annotation"

func Register() {

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
