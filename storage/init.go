package storage

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-storage/v2"
	"github.com/go-errors/errors"
	"github.com/samber/do/v2"
	"github.com/spf13/cast"
)

func Storage(name string) storage.FileStorage {
	return do.MustInvokeNamed[storage.FileStorage](global.Injector, "storage."+name)
}

func StorageRegister(name string, Type string, conf map[string]string) {
	do.OverrideNamed[storage.FileStorage](global.Injector, "storage."+name, func(injector do.Injector) (storage.FileStorage, error) {
		if conf == nil {
			return nil, errors.New(fmt.Sprintf("storage driver %s not found", name))
		}
		store, err := storage.New(Type, conf, LocalSign)
		return store, err
	})
}

func LocalSign(path string) (string, error) {
	data := map[string]any{
		"path":   path,
		"expire": time.Now().Add(3600 * time.Second).Unix(),
	}
	str, _ := sonic.Marshal(data)
	return helper.Encryption(string(str))
}

func LocalVerify(path string, sign string) bool {
	content, err := helper.Decryption(sign)
	if err != nil {
		return false
	}

	data := map[string]any{}
	_ = sonic.UnmarshalString(content, &data)

	if data["path"] != path {
		return false
	}

	if cast.ToInt64(data["expire"]) < time.Now().Unix() {
		return false
	}

	return true
}
