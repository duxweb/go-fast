package permission

import (
	"github.com/duxweb/go-fast/action"
	"github.com/duxweb/go-fast/i18n"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"sort"
)

type PermissionData struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
	Data  []*PermissionData
}

func New() *PermissionData {
	return &PermissionData{}
}

func (t *PermissionData) Group(name string, order int) *PermissionData {
	data := &PermissionData{
		Name:  name,
		Order: order,
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
	data := lo.Map[*PermissionData, map[string]any](t.Data, func(group *PermissionData, index int) map[string]any {
		list := lo.Map[*PermissionData, map[string]any](group.Data, func(item *PermissionData, index int) map[string]any {
			label := action.GetActionLabel(item.Name)
			return map[string]any{
				"name":  item.Name,
				"label": label,
			}
		})
		return map[string]any{
			"label":    i18n.Trans.Get(group.Name + ".name"),
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

func Can(permissions map[string]bool, name string) error {
	if len(permissions) == 0 {
		return nil
	}
	is, ok := permissions[name]
	if !ok {
		return nil
	}
	if !is {
		return echo.ErrForbidden
	}
	return nil
}
