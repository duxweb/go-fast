package config

import (
	"github.com/duxweb/go-fast/global"
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
	global.Debug = Load("app").GetBool("server.debug")
	global.DebugMsg = Load("app").GetString("server.debugMsg")

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
