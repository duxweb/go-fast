package dux

import (
	"embed"
	"github.com/duxweb/go-fast/app"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/service"
	"github.com/duxweb/go-fast/views"
	"github.com/duxweb/go-fast/web"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type Dux struct {
	registerApp []func()
	registerCmd []func(command *cobra.Command)
}

func New() *Dux {
	return &Dux{}
}

// RegisterApp Register Application
func (t *Dux) RegisterApp(calls ...func()) {
	t.registerApp = append(t.registerApp, calls...)
}

// RegisterCmd Register Command
func (t *Dux) RegisterCmd(calls ...func(command *cobra.Command)) {
	t.registerCmd = append(t.registerCmd, calls...)
}

// RegisterDir Register Directory
func (t *Dux) RegisterDir(dirs ...string) {
	app.DirList = append(app.DirList, dirs...)
}

//go:embed template/*
var FrameFs embed.FS

// Create Universal Service
func (t *Dux) create() {
	views.FrameFs = FrameFs
	for _, call := range t.registerApp {
		call()
	}
	t.RegisterCmd(app.Command, web.Command, database.Command)
}

// Run Command
func (t *Dux) Run() {
	t.create()
	var rootCmd = &cobra.Command{Use: "dux"}
	for _, cmd := range t.registerCmd {
		cmd(rootCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// SetTimezone Set Timezone
func (t *Dux) SetTimezone(location *time.Location) {
	time.Local = location
}

// SetTablePrefix Set Database Table Prefix
func (t *Dux) SetTablePrefix(prefix string) {
	global.TablePrefix = prefix
}

// SetConfigDir Set Configuration Directory
func (t *Dux) SetConfigDir(dir string) {
	global.ConfigDir = dir
}

// SetDatabaseStatus Set Database Status
func (t *Dux) SetDatabaseStatus(status bool) {
	service.Server.Database = status
}

// SetRedisStatus Set Redis Status
func (t *Dux) SetRedisStatus(status bool) {
	service.Server.Redis = status
}

// SetMongodbStatus Set MongoDB Status
func (t *Dux) SetMongodbStatus(status bool) {
	service.Server.Mongodb = status
}
