package config

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/duxweb/go-fast/v2/global"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"
)

var data = map[string]*koanf.Koanf{}

func Init() {
	configFiles, err := filepath.Glob(global.ConfigDir + "*.toml")
	if err != nil {
		panic("configuration loading failure")
	}

	for _, file := range configFiles {
		file = filepath.ToSlash(file)
		filename := path.Base(file)

		suffix := path.Ext(file)
		name := filename[0 : len(filename)-len(suffix)]
		LoadFile(name)
	}

	if IsLoad("use") {
		global.Debug = Load("use").Bool("app.debug")
		global.Lang = Load("use").String("app.lang")
	}
}

func LoadFile(name string) {
	// init config
	config := koanf.New(".")
	prefix := strings.ToUpper(name) + "_"

	// load toml file
	configFile := filepath.Join(global.ConfigDir, name+".toml")
	f := file.Provider(configFile)
	err := config.Load(f, toml.Parser())
	if err != nil {
		panic(err)
	}

	// load dotenv file
	envf := file.Provider(GetDotEnv())
	config.Load(envf, dotenv.ParserEnv(prefix, ".", func(s string) string {
		return strings.Replace(
			strings.ToLower(strings.TrimPrefix(s, prefix)),
			"_", ".", -1)
	}))

	// load env file
	envPrefix := "DUX_" + prefix
	config.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(
			strings.ToLower(strings.TrimPrefix(s, envPrefix)),
			"_", ".", -1)
	}), nil)

	f.Watch(func(event interface{}, err error) {
		if err != nil {
			slog.Error("config watch", slog.Any("error", err))
			return
		}
		config = koanf.New(".")
		config.Load(f, toml.Parser())
		data[name] = config
	})

	data[name] = config
}

func Load(name string) *koanf.Koanf {
	if t, ok := data[name]; ok {
		return t
	} else {
		panic("configuration (" + name + ") not found")
	}
}

func Reload(name string) {
	if _, ok := data[name]; ok {
		LoadFile(name)
	} else {
		panic("configuration (" + name + ") not found")
	}
}

func IsLoad(name string) bool {
	_, ok := data[name]
	return ok
}

func GetDotEnv() string {
	dotenv := filepath.Join(".", lo.Ternary(global.DotEnv == "", ".env", fmt.Sprintf(".%s.env", global.DotEnv)))
	if _, err := os.Stat(dotenv); os.IsNotExist(err) {
		return ""
	}
	return dotenv
}
