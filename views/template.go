package views

import (
	"embed"
	"encoding/json"
	"github.com/duxweb/go-fast/i18n"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"html/template"
	"io"
)

var Views = map[string]*template.Template{}

//go:embed template/*
var TplFs embed.FS

func Init() {
	NewFS("app", TplFs)
}

var funcMap = template.FuncMap{
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
func New(name string, dir string) *template.Template {
	if Views[name] == nil {
		Views[name] = template.Must(template.New("").Funcs(funcMap).ParseGlob(dir))
	}
	return Views[name]
}

// NewFS 创建虚拟模板
func NewFS(name string, fs embed.FS) *template.Template {
	if Views[name] == nil {
		Views[name] = template.Must(template.New("").Funcs(funcMap).ParseFS(fs, "**/*"))
	}
	return Views[name]
}

// Template 模板引擎
type Template struct{}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl := lo.Ternary[string](c.Get("tpl") != nil, cast.ToString(c.Get("tpl")), "default")
	if Views[tpl] == nil {
		return echo.ErrNotFound
	}
	return Views[tpl].ExecuteTemplate(w, name, data)
}

func Render() *Template {
	return &Template{}
}
