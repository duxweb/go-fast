package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxweb/go-fast/config"
	"github.com/samber/do"
	"runtime/debug"
)

func Init() {
	// Cache Size, Unit: M
	cacheSize := config.Get("app").GetInt("cache.size") * 1024 * 1024
	do.ProvideValue[*freecache.Cache](nil, freecache.NewCache(cacheSize))
	debug.SetGCPercent(20)
}
