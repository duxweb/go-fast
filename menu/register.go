package menu

var Menus = map[string]*MenuData{}

func Add(name string, route *MenuData) {
	Menus[name] = route
}

func Get(name string) *MenuData {
	return Menus[name]
}
