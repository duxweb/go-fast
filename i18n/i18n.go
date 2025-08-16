package i18n

import (
	"context"
	"embed"
	"io/fs"
	"sync"

	"github.com/duxweb/go-fast/v2/helper"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

var Localizer map[string]*i18n.Localizer

var LocalizerLock sync.RWMutex

type LangKey struct{}

//go:embed lang/*.toml
var langFs embed.FS

func Init() error {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Localizer = make(map[string]*i18n.Localizer)
	err := Register(langFs)
	return err
}

// Register 注册语言文件
// Register register language files
func Register(file embed.FS) error {
	err := fs.WalkDir(file, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		_, _ = Bundle.LoadMessageFileFS(file, path)
		return nil
	})
	return err
}

// WithLanguage 设置语言
// WithLanguage set language
func WithLanguage(ctx context.Context, lang string) context.Context {
	ctx = context.WithValue(ctx, LangKey{}, lang)
	return ctx
}

// GetLocalizer 获取本地化器
// GetLocalizer get localizer by lang
func GetLocalizer(lang string) *i18n.Localizer {
	LocalizerLock.RLock()
	defer LocalizerLock.RUnlock()
	localizer, ok := Localizer[lang]
	if !ok {
		localizer = i18n.NewLocalizer(Bundle, lang)
		Localizer[lang] = localizer
	}
	return localizer
}

// T 获取本地化消息
// T get localized message
func T(ctx context.Context, msg string, args ...any) string {

	lang := helper.GetContextValue[LangKey, string](ctx, LangKey{})
	if lang == "" {
		lang = "en-US"
	}

	defaultMessage := msg
	msgData := map[string]any{}

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]any:
			for k, v := range v {
				msgData[k] = v
			}
		case string:
			defaultMessage = v
		}
	}
	cfg := i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msg,
			Other: defaultMessage,
		},
		TemplateData: msgData,
	}
	str, err := GetLocalizer(lang).Localize(&cfg)

	if err != nil {
		return msg
	}
	return str
}
