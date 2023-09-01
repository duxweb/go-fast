package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxweb/go-fast/config"
	"runtime/debug"
)

var cache *freecache.Cache

func Init() {
	// Cache Size, Unit: M
	cacheSize := config.Load("app").GetInt("cache.size") * 1024 * 1024
	cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
}

func Injector() *freecache.Cache {
	return cache
}
