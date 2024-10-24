package views

import (
	"embed"
	"encoding/json"
	"github.com/duxweb/go-fast/i18n"
	"github.com/gofiber/template/jet/v2"
	"html/template"
	"net/http"
)

var Views = map[string]*jet.Engine{}

//go:embed template/*
var TplFs embed.FS

func Init() {
	NewFS("app", TplFs)
}

var funcMap = map[string]any{
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
	"t": func(s string) string {
		return i18n.Trans.Get(s)
	},
}

// New 创建普通模板
func New(name string, dir string) *jet.Engine {
	if Views[name] == nil {
		engine := jet.New(dir, ".jet")
		engine.AddFuncMap(funcMap)

		Views[name] = engine
	}
	return Views[name]
}

// NewFS 创建虚拟模板
func NewFS(name string, fs embed.FS) *jet.Engine {
	if Views[name] == nil {
		engine := jet.NewFileSystem(http.FS(fs), ".jet")
		engine.AddFuncMap(funcMap)
		Views[name] = engine
	}
	return Views[name]
}
