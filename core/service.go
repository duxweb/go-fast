package core

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/service"
)

func CoreService() *service.Config {
	return &service.Config{
		Name: "core",
		Init: func() error {
			// 自动创建目录
			// Automatically create directory
			for _, dir := range global.DirList {
				err := fileutil.CreateDir(dir)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
