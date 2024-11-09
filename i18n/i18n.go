package i18n

import (
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
	"io/fs"
)

var Bundle *i18n.Bundle

//go:embed lang/*.toml
var langFs embed.FS

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Register(langFs)
}

func Register(file embed.FS) {
	_ = fs.WalkDir(file, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		_, err = Bundle.LoadMessageFileFS(file, path)
		return nil
	})
}

func Get(c echo.Context, id string, data ...map[string]any) string {
	if c == nil {
		return id
	}
	return GetDefault(c, id, id, data...)
}

func GetDefault(c echo.Context, id string, message string, data ...map[string]any) string {
	local, ok := c.Get("i18n").(*i18n.Localizer)
	if !ok {
		return id
	}
	msgData := map[string]any{}
	if len(data) > 0 {
		msgData = data[0]
	}

	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			Other: message,
		},
		TemplateData: msgData,
	}
	str, err := local.Localize(cfg)
	if err != nil {
		return id
	}
	return str
}
