package config

import (
	"embed"
	"path"
	"path/filepath"

	"github.com/duxweb/go-fast/global"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/goutil/fsutil"
	"github.com/samber/do/v2"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var data = map[string]*viper.Viper{}

//go:embed all:tpl
var ConfigTplFs embed.FS

func Init() {
	// Init Config
	files, _ := ConfigTplFs.ReadDir("tpl")
	for _, file := range files {
		conf := filepath.Join(global.ConfigDir, file.Name())
		if fsutil.FileExist(conf) {
			continue
		}
		conf = filepath.ToSlash(conf)
		f := fsutil.MustCreateFile(conf, 0777, 0777)
		c, err := ConfigTplFs.ReadFile("tpl/" + file.Name())
		if err != nil {
			panic(err)
		}
		_, err = f.Write(c)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	configFiles, err := filepath.Glob(global.ConfigDir + "*.yaml")
	if err != nil {
		panic("configuration loading failure")
	}

	// Load Configuration Files in Loop
	for _, file := range configFiles {
		file = filepath.ToSlash(file)
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
	configFile := filepath.Join(global.ConfigDir, name+".yaml")
	config := viper.New()
	config.SetConfigFile(configFile)
	config.SetConfigType("yaml")
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
