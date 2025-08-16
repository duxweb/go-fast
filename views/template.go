package views

import (
	"embed"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/httpfs"
)

var Views = map[string]*jet.Set{}

//go:embed template/*
var TplFs embed.FS

func Init() error {
	NewFS("app", TplFs)
	return nil
}

// New 创建普通模板
func New(name string, dir string) *jet.Set {
	if Views[name] == nil {
		loader := jet.NewOSFileSystemLoader(dir)
		engine := jet.NewSet(loader)
		Views[name] = engine
	}
	return Views[name]
}

// NewFS 创建虚拟模板
func NewFS(name string, fs embed.FS) *jet.Set {
	if Views[name] == nil {
		loader, err := httpfs.NewLoader(http.FS(fs))
		if err != nil {
			panic(err)
		}
		engine := jet.NewSet(loader)
		Views[name] = engine
	}
	return Views[name]
}
