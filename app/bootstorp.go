package app

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/duxweb/go-fast/annotation"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/permission"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/web"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

// Dux 基础结构
// Dux infrastructure
type Dux struct {
	apps []func()
	cmds []func() []*cli.Command
	Lang string
}

// RegisterApp 应用注册
// RegisterApp register application
func (t *Dux) RegisterApp(calls ...func()) {
	t.apps = append(t.apps, calls...)
}

// RegisterCmd 命令注册
// RegisterCmd register command
func (t *Dux) RegisterCmd(calls ...func() []*cli.Command) {
	t.cmds = append(t.cmds, calls...)
}

// RegisterDir 自动创建目录
// RegisterDir register folder creation
func (t *Dux) RegisterDir(dirs ...string) {
	helper.CreateDir(dirs...)
}

// RegisterTpl 注册模板目录
// RegisterTpl register template folder
func (t *Dux) RegisterTpl(name string, dir string) {
	views.New(name, dir)
}

// RegisterTplFS 注册虚拟模板目录
// RegisterTplFS register FS template folder
func (t *Dux) RegisterTplFS(name string, fs embed.FS) {
	views.NewFS(name, fs)
}

// RegisterLangFS 注册语言包
// RegisterLangFS register language packs
func (t *Dux) RegisterLangFS(fs embed.FS) {
	i18n.Register(fs)
}

// 创建命令
// create command
func (t *Dux) create() {
	for _, call := range t.apps {
		call()
	}
	t.RegisterCmd(Command, web.Command, permission.Command, database.Command)
}

// Run 运行命令
// Run command
func (t *Dux) Run() {

	t.create()

	list := make([]*cli.Command, 0)
	for _, cmd := range t.cmds {
		list = append(list, cmd()...)
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

// SetAnnotations 设置索引文件
// SetAnnotations set index file
func (t *Dux) SetAnnotations(data []*annotation.File) {
	annotation.Annotations = data
}

// SetStaticFs 设置静态目录
// SetStaticFs setting up static folder
func (t *Dux) SetStaticFs(fs embed.FS) {
	global.StaticFs = &fs
}

// SetPageFs 设置页面目录
// SetPageFs setting page folder
func (t *Dux) SetPageFs(fs embed.FS) {
	global.PageFs = &fs
}

// SetTimezone 设置时区
// SetStaticFs setting the time zone
func (t *Dux) SetTimezone(location *time.Location) {
	time.Local = location
}

// SetTablePrefix 设置表前缀
// SetTablePrefix Set table prefix
func (t *Dux) SetTablePrefix(prefix string) {
	global.TablePrefix = prefix
}

// SetConfigDir 设置配置目录
// SetConfigDir setting the config folder
func (t *Dux) SetConfigDir(dir string) {
	global.ConfigDir = dir
}

// SetDataDir 设置数据目录
// SetConfigDir setting the data folder
func (t *Dux) SetDataDir(dir string) {
	global.DataDir = dir
}

// IsRelease 判断是否编译
// IsRelease determine whether to build
func IsRelease() bool {
	arg1 := strings.ToLower(os.Args[0])
	name := filepath.Base(arg1)
	return strings.Index(name, "__") != 0 && strings.Index(arg1, "go-build") < 0
}
