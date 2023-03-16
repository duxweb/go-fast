package permission

var Permissions = map[string]*PermissionData{}

func Add(name string, route *PermissionData) {
	Permissions[name] = route
}

func Get(name string) *PermissionData {
	return Permissions[name]
}
