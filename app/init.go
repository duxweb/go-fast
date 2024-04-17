package app

import (
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/event"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/resources"
	"github.com/duxweb/go-fast/route"
	"github.com/duxweb/go-fast/task"
)

var DirList = []string{
	"./uploads",
	"./data",
	"./config",
	"./app",
	"./tmp",
	"./data/logs",
	"./data/logs/default",
	"./data/logs/request",
	"./data/logs/service",
	"./data/logs/database",
	"./data/logs/task"}

func Init() {

	// 自动创建目录
	for _, path := range global.DirList {
		if !helper.IsExist(path) {
			if !helper.CreateDir(path) {
				panic("failed to create " + path + " directory")
			}
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
