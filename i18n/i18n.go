package i18n

import (
	"embed"
	"github.com/duxweb/go-fast/global"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/fs"
)

var Bundle *i18n.Bundle

var Trans *Localizer

//go:embed lang/*.yaml
var langFs embed.FS

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	Trans = &Localizer{bundle: Bundle, localizer: i18n.NewLocalizer(Bundle, global.Lang)}
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

type Localizer struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

func (l Localizer) Get(id string) string {
	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			One:   id,
			Other: id,
		},
	}
	str, err := l.localizer.Localize(cfg)
	if err != nil {
		return id
	}

	return str
}

func (l Localizer) GetData(id string, data map[string]any) string {
	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			One:   id,
			Other: id,
		},
		TemplateData: data,
	}
	str, err := l.localizer.Localize(cfg)
	if err != nil {
		return id
	}

	return str
}
