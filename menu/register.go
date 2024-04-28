package menu

var Menus = map[string]*MenuData{}

func Set(name string, data *MenuData) {
	Menus[name] = data
}

func Get(name string) *MenuData {
	return Menus[name]
}
