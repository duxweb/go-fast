package route

var Routes = map[string]*RouterData{}

func Set(name string, route *RouterData) *RouterData {
	Routes[name] = route
	return route
}

func Get(name string) *RouterData {
	return Routes[name]
}
