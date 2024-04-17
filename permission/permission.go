package permission

import (
	"github.com/samber/lo"
	"sort"
)

type PermissionData struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Order int    `json:"order"`
	Data  []*PermissionData
}

func New() *PermissionData {
	return &PermissionData{}
}

func (t *PermissionData) Group(name string, order int) *PermissionData {
	data := &PermissionData{
		Name:  name,
		Label: "",
		Order: order,
	}
	t.Data = append(t.Data, data)
	return data
}

func (t *PermissionData) Add(name string) {
	data := &PermissionData{
		Name:  t.Name + "." + name,
		Label: "",
	}
	t.Data = append(t.Data, data)
}
func (t *PermissionData) Get() []map[string]any {
	data := lo.Map[*PermissionData, map[string]any](t.Data, func(group *PermissionData, index int) map[string]any {
		list := lo.Map[*PermissionData, map[string]any](group.Data, func(item *PermissionData, index int) map[string]any {
			return map[string]any{
				"name":  item.Name,
				"label": item.Label,
			}
		})
		return map[string]any{
			"label":    group.Label,
			"order":    group.Order,
			"name":     "group:" + group.Name,
			"children": list,
		}
	})
	sort.Slice(data, func(i, j int) bool {
		return data[i]["order"].(int) < data[j]["order"].(int)
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
