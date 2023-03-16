package app

import (
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
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

	// Automatically Create Directory
	for _, path := range global.DirList {
		if !helper.IsExist(path) {
			if !helper.CreateDir(path) {
				panic("failed to create " + path + " directory")
			}
		}
	}

	// Initialize the application and register the routing and other initialization processes for the application in this closure.
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Init != nil {
			appConfig.Init()
		}
	}

	// Register the application and register application routes and other data in this closure.
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Register != nil {
			appConfig.Register()
		}
	}

	// Start the application, which includes the startup process after the application is registered, and handle post-calling methods.
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Boot != nil {
			appConfig.Boot()
		}
	}

}
