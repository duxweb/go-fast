package dux

import (
	"embed"
	"github.com/duxweb/go-fast/annotation"
	"github.com/duxweb/go-fast/app"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/service"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/web"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed template/*
var TplFs embed.FS

//go:embed lang/*.yaml
var LangFs embed.FS

// Dux 基础结构
type Dux struct {
	apps []func()
	cmds []func() []*cli.Command
}

func New() *Dux {
	return &Dux{}
}

// RegisterApp 注册模块
func (t *Dux) RegisterApp(calls ...func()) {
	t.apps = append(t.apps, calls...)
}

// RegisterCmd 注册命令
func (t *Dux) RegisterCmd(calls ...func() []*cli.Command) {
	t.cmds = append(t.cmds, calls...)
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

// RegisterLangFS 注册语言包
func (t *Dux) RegisterLangFS(fs embed.FS) {
	i18n.FsList = append(i18n.FsList, fs)
}

// RegisterAnnotations 注册索引
func (t *Dux) RegisterAnnotations(data []*annotation.File) {
	annotation.Annotations = data
}

// 创建公共服务
func (t *Dux) create() {
	t.RegisterTplFS("app", TplFs)
	t.RegisterLangFS(LangFs)
	for _, call := range t.apps {
		call()
	}
	t.RegisterCmd(web.Command, app.Command, annotation.Command, database.Command)
}

// Run Command
func (t *Dux) Run() {
	t.create()

	list := make([]*cli.Command, 0)
	for _, cmd := range t.cmds {
		list = append(list, cmd()...)
	}

	if !IsRelease() {
		annotation.Run()
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

func IsRelease() bool {
	arg1 := strings.ToLower(os.Args[0])
	name := filepath.Base(arg1)
	return strings.Index(name, "__") != 0 && strings.Index(arg1, "go-build") < 0
}
