package permission

var Permissions = map[string]*PermissionData{}

func Set(name string, data *PermissionData) {
	Permissions[name] = data
}

func Get(name string) *PermissionData {
	return Permissions[name]
}
