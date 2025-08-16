package route

import (
	"github.com/duxweb/go-fast/v2/annotation"
)

var Routes = map[string]*RouterData{}

func SetRouter(name string, route *RouterData) *RouterData {
	Routes[name] = route
	return route
}

func GetRouter(name string) *RouterData {
	return Routes[name]
}

func Register() {
	for _, file := range annotation.Annotations {
		for _, item := range file.Annotations {
			if item.Name != "Route" {
				continue
			}
			fun, ok := item.Func.(func())
			if !ok {
				continue
			}
			fun()
		}
	}
}
