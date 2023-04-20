package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/duxweb/go-fast/global"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Init() {
	bundle := i18n.NewBundle(global.Lang)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Loading Multilingual
	bundle.MustParseMessageFileBytes([]byte(`
		HelloWorld = "Hello World!"
		`), "en.toml")
}
