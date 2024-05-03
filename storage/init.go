package storage

import (
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-storage"
	"github.com/go-errors/errors"
	"github.com/samber/do/v2"
	"github.com/spf13/cast"
)

func Init() {
	dbConfig := config.Load("storage").GetStringMap("drivers")
	for name, conf := range dbConfig {
		do.ProvideNamed[storage.FileStorage](global.Injector, "storage."+name, func(injector do.Injector) (storage.FileStorage, error) {
			if conf == nil {
				return nil, errors.New(fmt.Sprintf("storage driver %s not found", name))
			}
			return storage.New(name, cast.ToStringMapString(conf)), nil
		})
	}
}

func Storage(names ...string) storage.FileStorage {
	name := config.Load("storage").GetString("type")
	if len(names) > 0 {
		name = names[0]
	}
	return do.MustInvokeNamed[storage.FileStorage](global.Injector, "storage."+name)
}
