package app

import (
	"embed"
	"github.com/duxweb/go-fast/annotation"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/permission"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/web"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Dux 基础结构
type Dux struct {
	apps []func()
	cmds []func() []*cli.Command
	Lang string
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
	CreateDir(dirs...)
}

// RegisterTpl 注册模板目录
func (t *Dux) RegisterTpl(name string, dir string) {
	views.New(name, dir)
}

// RegisterTplFS 注册虚拟模板目录
func (t *Dux) RegisterTplFS(name string, fs embed.FS) {
	views.NewFS(name, fs)
}

// RegisterLangFS 注册语言包
func (t *Dux) RegisterLangFS(fs embed.FS) {
	i18n.Register(fs)
}

// 创建公共服务
func (t *Dux) create() {
	for _, call := range t.apps {
		call()
	}
	t.RegisterCmd(Command, web.Command, permission.Command, annotation.Command, database.Command)
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

		if global.AnnotationUpdate {
			color.Redln("⇨ runtime found an update, please restart")
			os.Exit(0)
		}
	}

	Start(t)

	appCli := &cli.App{
		Name:     "dux",
		Commands: list,
	}
	if err := appCli.Run(os.Args); err != nil {
		color.Errorln(err)
	}
}

// SetAnnotations 设置索引
func (t *Dux) SetAnnotations(data []*annotation.File) {
	annotation.Annotations = data
}

// SetStaticFs 设置静态目录
func (t *Dux) SetStaticFs(fs embed.FS) {
	global.StaticFs = &fs
}

// SetTimezone 设置时区
func (t *Dux) SetTimezone(location *time.Location) {
	time.Local = location
}

// SetTablePrefix 设置表前缀
func (t *Dux) SetTablePrefix(prefix string) {
	global.TablePrefix = prefix
}

// SetConfigDir 设置配置目录
func (t *Dux) SetConfigDir(dir string) {
	global.ConfigDir = dir
}

// SetDataDir 设置数据目录
func (t *Dux) SetDataDir(dir string) {
	global.DataDir = dir
}

func IsRelease() bool {
	arg1 := strings.ToLower(os.Args[0])
	name := filepath.Base(arg1)
	return strings.Index(name, "__") != 0 && strings.Index(arg1, "go-build") < 0
}