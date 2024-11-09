package web

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/views"
	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"io"
	"sync"
)

func ViewHandler() *Template {
	return &Template{
		Mutex: sync.RWMutex{},
	}
}

// Template 模板引擎
type Template struct {
	Mutex sync.RWMutex
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl := lo.Ternary[string](c.Get("tpl") != nil, cast.ToString(c.Get("tpl")), "default")
	if views.Views[tpl] == nil {
		return errors.New("tpl app not found")
	}
	t.Mutex.RLock()
	template, err := views.Views[tpl].GetTemplate(name)
	views.Views[tpl].AddGlobal("t", func(s string) string {
		return i18n.Get(c, s)
	})
	t.Mutex.RUnlock()
	if err != nil {
		return err
	}

	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	bind := jetVarMap(data)

	return template.Execute(w, bind, nil)
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
	case echo.Map:
		bind = make(jet.VarMap)
		for key, value := range binds {
			bind.Set(key, value)
		}
	case jet.VarMap:
		bind = binds
	}
	return bind
}
