package web

import (
	"io"
	"strings"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/duxweb/go-fast/v2/i18n"
	"github.com/duxweb/go-fast/v2/views"
	"github.com/labstack/echo/v4"
)

type Render struct {
	Mutex sync.RWMutex
}

func ViewHandler() *Render {
	return &Render{
		Mutex: sync.RWMutex{},
	}
}

func (t *Render) Render(w io.Writer, name string, data interface{}, r echo.Context) error {

	fileName := strings.Split(name, ":")

	tplName := "app"
	path := ""
	if len(fileName) > 1 {
		tplName = fileName[0]
		path = fileName[1]
	} else {
		path = fileName[0]
	}

	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	template, err := views.Views[tplName].GetTemplate(path)
	views.Views[tplName].AddGlobal("t", func(s string) string {
		return i18n.T(r.Request().Context(), s)
	})
	if err != nil {
		return err
	}

	return template.Execute(w, jetVarMap(data), nil)
}

func jetVarMap(binding any) jet.VarMap {
	var bind jet.VarMap
	if binding == nil {
		return bind
	}
	switch binds := binding.(type) {
	case map[string]any:
		bind = make(jet.VarMap)
		for key, value := range binds {
			bind.Set(key, value)
		}
	case jet.VarMap:
		bind = binds
	}
	return bind
}
