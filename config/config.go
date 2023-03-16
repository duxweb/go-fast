package config

import (
	"fmt"
	"github.com/duxweb/go-fast/global"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

type Config map[string]*viper.Viper

func Init() {
	do.ProvideValue[Config](nil, map[string]*viper.Viper{})

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
		do.MustInvoke[Config](nil)[name] = LoadConfig(name)
	}

	// Set Framework Configuration
	global.Debug = Get("app").GetBool("server.debug")
	global.DebugMsg = Get("app").GetString("server.debugMsg")

}

// LoadConfig Load Configuration from Specified File
func LoadConfig(name string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(name)
	config.SetConfigType("yaml")
	config.AddConfigPath(global.ConfigDir)
	if err := config.ReadInConfig(); err != nil {
		fmt.Println("config", name)
		panic(err)
	}
	return config
}

// Get File Configuration
func Get(name string) *viper.Viper {
	if t, ok := do.MustInvoke[Config](nil)[name]; ok {
		return t
	} else {
		panic("configuration (" + name + ") not found")
	}
}
