package config

import (
	"github.com/duxweb/go-fast/global"
	"github.com/golang-module/carbon/v2"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

var data = map[string]*viper.Viper{}

func Init() {

	pwd, _ := os.Getwd()
	configFiles, err := filepath.Glob(filepath.Join(pwd, global.ConfigDir+"*.yaml"))
	if err != nil {
		panic("configuration loading failure")
	}

	// Load Configuration Files in Loop
	for _, file := range configFiles {
		filename := path.Base(file)
		suffix := path.Ext(file)
		name := filename[0 : len(filename)-len(suffix)]
		data[name] = LoadFile(name)
	}

	// Set Framework Configuration
	if IsLoad("use") {
		global.Debug = Load("use").GetBool("app.debug")
		global.Lang = Load("use").GetString("app.lang")
	}

	// Set time
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.DateTimeLayout,
		Timezone:     carbon.Local,
		WeekStartsAt: carbon.Monday,
		Locale:       global.Lang,
	})

	// set OsFs
	do.ProvideNamed[afero.Fs](global.Injector, "os.fs", func(injector do.Injector) (afero.Fs, error) {
		fs := afero.NewOsFs()
		return fs, nil
	})
}

func LoadFile(name string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(name)
	config.SetConfigType("yaml")
	config.AddConfigPath(global.ConfigDir)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	return config
}

func Load(name string) *viper.Viper {
	if t, ok := data[name]; ok {
		return t
	} else {
		panic("configuration (" + name + ") not found")
	}
}

func IsLoad(name string) bool {
	_, ok := data[name]
	return ok
}
