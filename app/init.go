package app

import (
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/event"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/resources"
	"github.com/duxweb/go-fast/route"
	"github.com/duxweb/go-fast/task"
	"github.com/samber/do"
	"github.com/spf13/afero"
)

var DirList = []string{
	"./app",
	"./public",
	"./public/uploads",
	"./data",
	"./data/tmp",
	"./config",
	"./data/logs",
}

func Init() {

	// 自动创建目录
	fs := do.MustInvokeNamed[afero.Fs](global.Injector, "os.fs")
	for _, path := range global.DirList {
		exists, _ := afero.DirExists(fs, path)
		if exists {
			continue
		}
		err := fs.MkdirAll(path, 0777)
		if err != nil {
			panic("failed to create " + path + " directory")
		}
	}

	// 初始化应用模块
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Init != nil {
			appConfig.Init()
		}
	}

	// 注册应用模块
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Register != nil {
			appConfig.Register()
		}
	}

	// 启动应用模块
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Boot != nil {
			appConfig.Boot()
		}
	}

	database.Register()
	event.Register()
	task.Register()
	resources.Register()
	route.Register()

}
