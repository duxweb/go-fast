package app

import (
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/event"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/resources"
	"github.com/duxweb/go-fast/route"
	"github.com/duxweb/go-fast/task"
)

var dirList = []string{
	"./database",
	"./public",
	"./public/uploads",
	"./data",
	"./data/tmp",
	"./config",
	"./data/logs",
}

func Init(t *Dux) {

	// 自动创建目录
	helper.CreateDir(dirList...)

	// 初始化应用模块
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Init != nil {
			appConfig.Init(t)
		}
	}

	// 注册应用模块
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Register != nil {
			appConfig.Register(t)
		}
	}

	// 启动应用模块
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Boot != nil {
			appConfig.Boot(t)
		}
	}

	database.Register()
	event.Register()
	task.Register()
	resources.Register()
	route.Register()

}
