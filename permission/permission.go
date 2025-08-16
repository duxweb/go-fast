package permission

import (
	"github.com/samber/lo"
)

type PermissionData struct {
	Name string `json:"name"`
	Data []*PermissionData
}

func New() *PermissionData {
	return &PermissionData{}
}

func (t *PermissionData) Group(name string, order int) *PermissionData {
	data := &PermissionData{
		Name: name,
	}
	t.Data = append(t.Data, data)
	return data
}

func (t *PermissionData) Add(name string) {
	allName := t.Name + "." + name
	data := &PermissionData{
		Name: allName,
	}
	t.Data = append(t.Data, data)
}
func (t *PermissionData) Get() []map[string]any {
	data := lo.Map(t.Data, func(group *PermissionData, index int) map[string]any {
		list := lo.Map(group.Data, func(item *PermissionData, index int) map[string]any {
			return map[string]any{
				"name": item.Name,
			}
		})

		return map[string]any{
			"name":     "group:" + group.Name,
			"children": list,
		}
	})
	return data
}

func (t *PermissionData) GetData() []string {
	data := make([]string, 0)
	for _, datum := range t.Data {
		for _, item := range datum.Data {
			data = append(data, item.Name)
		}
	}
	return data
}
