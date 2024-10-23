package views

import (
	"embed"
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/duxweb/go-fast/i18n"
	"github.com/gofiber/fiber/v2"
	core "github.com/gofiber/template"
	"io"
)

var Views = map[string]*jet.Set{}

//go:embed template/*
var TplFs embed.FS

func Init() {
	NewFS("app", TplFs)
}

// New 创建普通模板
func New(name string, dir string) *jet.Set {
	if Views[name] == nil {
		view := jet.NewHTMLSet(dir)
		view.AddGlobal("t", func(s string) string {
			return i18n.Trans.Get(s)
		})
		Views[name] = view
	}
	return Views[name]
}

// NewFS 创建虚拟模板
func NewFS(name string, fs embed.FS) *jet.Set {
	if Views[name] == nil {
		view := jet.NewHTMLSetLoader(&Loader{
			fs: fs,
		})
		view.AddGlobal("t", func(s string) string {
			return i18n.Trans.Get(s)
		})
		Views[name] = view
	}
	return Views[name]
}

// Loader 加载引擎
type Loader struct {
	fs embed.FS
}

func (l *Loader) Open(name string) (io.ReadCloser, error) {
	file, err := l.fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, err
}

func (l *Loader) Exists(name string) (string, bool) {
	_, err := l.fs.Open(name)
	if err != nil {
		return name, false
	}
	return name, true
}

// Engine 模板引擎
type Engine struct {
	core.Engine
}

func (t *Engine) Load() error {
	return nil
}

func (t *Engine) Render(out io.Writer, name string, binding interface{}, apps ...string) error {

	app := "app"
	if len(apps) > 0 {
		app = apps[0]
	}

	t.Mutex.RLock()
	tmpl, err := Views[app].GetTemplate(name)
	t.Mutex.RUnlock()

	if err != nil || tmpl == nil {
		return fmt.Errorf("render: template %s could not be Loaded: %w", name, err)
	}

	bind := jetVarMap(binding)

	return tmpl.Execute(out, bind, nil)
}

func Fiber() fiber.Views {
	return &Engine{}
}

func jetVarMap(binding interface{}) jet.VarMap {
	var bind jet.VarMap
	if binding == nil {
		return bind
	}
	switch binds := binding.(type) {
	case map[string]interface{}:
		bind = make(jet.VarMap)
		for key, value := range binds {
			bind.Set(key, value)
		}
	case fiber.Map:
		bind = make(jet.VarMap)
		for key, value := range binds {
			bind.Set(key, value)
		}
	case jet.VarMap:
		bind = binds
	}
	return bind
}
