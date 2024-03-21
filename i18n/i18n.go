package i18n

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/duxweb/go-fast/global"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locale/*.toml
var LocaleFS embed.FS

var Trans *Localizer

type Localizer struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

func Init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, "common.en-US.toml")
	bundle.LoadMessageFileFS(LocaleFS, "common.zh-CN.toml")
	Trans = &Localizer{bundle: bundle, localizer: i18n.NewLocalizer(bundle, global.Lang)}
}

func (l Localizer) Get(id string) string {
	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			Other: id,
			One:   id,
		},
	}
	str, err := l.localizer.Localize(cfg)
	if err != nil {
		return id
	}

	return str
}

func (l Localizer) GetWithData(id string, data map[string]any) string {
	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
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
