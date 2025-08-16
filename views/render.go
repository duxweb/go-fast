package views

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/duxweb/go-fast/v2/i18n"
)

var View = ViewHandler()

// HTMLOutput HTML 响应输出结构
type HTMLOutput struct {
	ContentType string `header:"Content-Type"`
	Body        string `doc:"HTML content"`
}

// NewHTMLResponse 创建 HTML 响应
func NewHTMLResponse(html string) *HTMLOutput {
	return &HTMLOutput{
		ContentType: "text/html; charset=utf-8",
		Body:        html,
	}
}

// Render 使用 views 系统渲染模板并返回 HTML 响应
func Render(ctx context.Context, templateName string, data any) (*HTMLOutput, error) {
	var buf bytes.Buffer

	// 使用现有的 View.Render 方法
	err := View.Render(&buf, templateName, data, ctx)
	if err != nil {
		return nil, err
	}

	return NewHTMLResponse(buf.String()), nil
}

type RenderHandler struct {
	Mutex sync.RWMutex
}

func ViewHandler() *RenderHandler {
	return &RenderHandler{
		Mutex: sync.RWMutex{},
	}
}

func (t *RenderHandler) Render(w io.Writer, name string, data any, c context.Context) error {

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
	template, err := Views[tplName].GetTemplate(path)
	Views[tplName].AddGlobal("t", func(s string) string {
		return i18n.T(c, s)
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
