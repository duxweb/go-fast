package dux

import (
	"embed"
	"github.com/duxweb/go-fast/app"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/service"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/web"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

//go:embed template/*
var TplFs embed.FS

// Dux 基础结构
type Dux struct {
	registerApp []func()
	registerCmd []func() []*cli.Command
}

func New() *Dux {
	return &Dux{}
}

// RegisterApp 注册模块
func (t *Dux) RegisterApp(calls ...func()) {
	t.registerApp = append(t.registerApp, calls...)
}

// RegisterCmd 注册命令
func (t *Dux) RegisterCmd(calls ...func() []*cli.Command) {
	t.registerCmd = append(t.registerCmd, calls...)
}

// RegisterDir 注册目录
func (t *Dux) RegisterDir(dirs ...string) {
	app.DirList = append(app.DirList, dirs...)
}

// RegisterTpl 注册模板目录
func (t *Dux) RegisterTpl(name string, dir string) {
	views.DirList[name] = dir
}

// RegisterTplFS 注册虚拟模板目录
func (t *Dux) RegisterTplFS(name string, fs embed.FS) {
	views.FsList[name] = fs
}

// 创建公共服务
func (t *Dux) create() {
	t.RegisterTplFS("app", TplFs)
	for _, call := range t.registerApp {
		call()
	}
	t.RegisterCmd(app.Command, web.Command, database.Command)
}

// Run Command
func (t *Dux) Run() {
	t.create()

	list := make([]*cli.Command, 0)
	for _, cmd := range t.registerCmd {
		list = append(list, cmd()...)
	}
	appCli := &cli.App{
		Name:     "dux",
		Commands: list,
	}
	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// SetTimezone 设置时区
func (t *Dux) SetTimezone(location *time.Location) {
	time.Local = location
}

// SetLang 设置语言
func (t *Dux) SetLang(lang string) {
	global.Lang = lang
}

// SetTablePrefix 设置表前缀
func (t *Dux) SetTablePrefix(prefix string) {
	global.TablePrefix = prefix
}

// SetConfigDir 设置配置目录
func (t *Dux) SetConfigDir(dir string) {
	global.ConfigDir = dir
}

// SetDatabaseStatus 设置数据库开关
func (t *Dux) SetDatabaseStatus(status bool) {
	service.Server.Database = status
}

// SetRedisStatus 设置redis开关
func (t *Dux) SetRedisStatus(status bool) {
	service.Server.Redis = status
}

// SetMongodbStatus 设置mongodb开关
func (t *Dux) SetMongodbStatus(status bool) {
	service.Server.Mongodb = status
}
