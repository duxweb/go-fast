<p align="center">
<a href="https://www.duxravel.com/">
    <img src="https://github.com/duxphp/duxravel/blob/main/resources/image/watermark.png?raw=true" width="100" height="100">
</a>

<p align="center">
  <img alt="Version" src="https://img.shields.io/badge/version-lpha-red.svg?cacheSeconds=2592000" />
  <a href="https://github.com/duxweb/go-storage/blob/main/LICENSE" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
    <a title="Go Reference" target="_blank" href="https://pkg.go.dev"><img src="https://img.shields.io/github/go-mod/go-version/duxweb/go-storage"></a>
</p>

<p align="center"><code>DuxFast</code> 是一款基于 GoEcho 的快速开发框架，集成主流三方包，简单、易开发、高性能的集成框架。</p>



# 💥 版本

警告：该版本作为开发版，尚有功能正在开发中并有不可避免的 bug，请勿在正式环境中使用。

# 🎯 特点

- 📦 基于 GoEcho 的高性能 Web 框架。
- 📚 整合 Gorm 作为主要数据库驱动，提供良好的数据库操作支持。
- 🧭 自动注解，自动生成注解索引，用于路由、任务、事件等自动化处理。
- 📡 不做过度封装，便于开发者灵活选择和随版本升级。
- 🔧 集成各大流行包，并封装常用日志、异常、权限等工具包。
- 📡 采用应用模块化设计，提高应用程序的可维护性和可扩展性。
- 📡 统一注册应用入口，方便应用程序的整体架构和管理。
- 🏷 开发命令助手与脚手架工具，提供基础的代码生成。


#  ⚡ 快速开始

```go
package main

import (
	"github.com/duxweb/go-fast/app"
	"project/app/home"
)

func main() {
	dux := duxgo.New()
	dux.RegisterApp(home.App)
	dux.Run()
}

```


```go
package home

import (
	"github.com/duxweb/go-fast/app"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
)

var config = struct {
}{}

func App() {
	app.Register(&app.Config{
		Name:     "home",
		Title:    "Example",
		Desc:     "This is an example",
		Config:   &config,
		Init:     Init,
		Register: Register,
	})
}

func Init() {
	route.Add("web", route.New(""))
}

func Register() {
	group := route.Get("web")
	group.Get("/", func(c *echo.Context) error {
		return c.String(200, "I'm a GET request!")
	}, "web.home")

}

```

#  ⚙ 安装

请确保当前 Golang 环境版本高于 `1.18` 版本，建立项目目录并初始化。

```sh
go get github.com/duxweb/go-fast
```

# 💡思想

该框架遵循与 DuxLite 一致化架构设计，将各个功能模块应用化，并通过 `应用入口` 与 `事件调度` 进行高度解耦，并保证基础框架与系统必备最小化，避免大而全的臃肿框架设计。