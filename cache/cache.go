package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxweb/go-fast/config"
	"runtime/debug"
)

var cache *freecache.Cache

func Init() {
	// Cache Size, Unit: M
	size := 100
	if config.IsLoad("cache") {
		size = config.Load("cache").GetInt("size")
	}
	cacheSize := size * 1024 * 1024
	cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
}

func Injector() *freecache.Cache {
	return cache
}
