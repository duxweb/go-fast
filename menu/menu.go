package menu

import (
	"sort"
)

// MenuData Menu application structure.
type MenuData struct {
	App      string `json:"app"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Icon     string `json:"icon"`
	Title    string `json:"title"`
	Hidden   bool   `json:"hidden"`
	Order    int    `json:"order"`
	Data     []*MenuData
	PushData map[string]*MenuData
}

// New Create a new menu.
func New() *MenuData {
	return &MenuData{
		PushData: map[string]*MenuData{},
	}
}

func (t *MenuData) Add(data *MenuData) *MenuData {
	t.Data = append(t.Data, data)
	return data
}

func (t *MenuData) Push(app string) *MenuData {
	data := &MenuData{App: app}
	t.PushData[data.App] = data
	return data
}

func (t *MenuData) Group(name string) *MenuData {
	data := &MenuData{
		Name: name,
	}
	t.Data = append(t.Data, data)
	return data
}

func (t *MenuData) Item(name string, url string, order int) {
	data := &MenuData{
		Name:  name,
		Url:   url,
		Order: order,
	}
	t.Data = append(t.Data, data)
}

func (t *MenuData) Get() []map[string]any {
	// Reset the menu
	var menu []map[string]any
	for _, appData := range t.Data {
		// Merge and append menus
		if t.PushData[appData.App] != nil {

			for _, datum := range t.PushData[appData.App].Data {
				appData.Data = append(appData.Data, datum)
			}
		}
		// Reset group menus
		var group []map[string]any
		for _, groupData := range appData.Data {
			// Reset submenu
			var list []map[string]any
			for _, items := range groupData.Data {
				list = append(list, map[string]any{
					"name":  items.Name,
					"url":   items.Url,
					"order": items.Order,
				})
			}
			sort.Slice(list, func(i, j int) bool {
				return list[i]["order"].(int) < list[j]["order"].(int)
			})
			group = append(group, map[string]any{
				"name":  groupData.Name,
				"order": groupData.Order,
				"title": groupData.Title,
				"menu":  list,
			})
		}
		sort.Slice(group, func(i, j int) bool {
			return group[i]["order"].(int) < group[j]["order"].(int)
		})
		menu = append(menu, map[string]any{
			"name":  appData.Name,
			"icon":  appData.Icon,
			"order": appData.Order,
			"url":   appData.Url,
			"menu":  group,
		})
	}
	sort.Slice(menu, func(i, j int) bool {
		return menu[i]["order"].(int) < menu[j]["order"].(int)
	})

	return menu
}
