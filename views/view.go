package views

import (
	"embed"
	"encoding/json"
	"github.com/gofiber/template/html"
	"github.com/yalue/merged_fs"
	"html/template"
	"net/http"
)

var TplFs embed.FS

var FrameFs embed.FS

var Views *html.Engine

func Init() {
	// Registration fiber template
	mergedFS := merged_fs.NewMergedFS(FrameFs, TplFs)

	engine := html.NewFileSystem(http.FS(mergedFS), ".gohtml")

	engine.AddFunc("unescape", func(v string) template.HTML {
		return template.HTML(v)
	})
	engine.AddFunc("marshal", func(v string) string {
		a, _ := json.Marshal(v)
		return string(a)
	})
	Views = engine
}
