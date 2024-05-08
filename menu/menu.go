package menu

import (
	"github.com/samber/lo"
	"sort"
	"strings"
)

// MenuData Menu application structure.
type MenuData struct {
	Prefix   string               `json:"-"`
	App      string               `json:"app"`
	Name     string               `json:"name"`
	Label    string               `json:"label"`
	Route    string               `json:"route"`
	Icon     string               `json:"icon"`
	Sort     int                  `json:"sort"`
	Meta     map[string]any       `json:"meta"`
	Data     []*MenuData          `json:"-"`
	PushData map[string]*MenuData `json:"-"`
}

// New Create a new menu.
func New(prefix string) *MenuData {
	return &MenuData{
		PushData: map[string]*MenuData{},
		Prefix:   prefix,
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

func (t *MenuData) Group(name string, label string, icon string) *MenuData {
	data := &MenuData{
		Name:  name,
		Label: label,
		Icon:  icon,
	}
	t.Data = append(t.Data, data)
	return data
}

func (t *MenuData) Item(name string, label string, route string, sort int) {
	data := &MenuData{
		Name:  name,
		Route: route,
		Label: label,
		Sort:  sort,
	}
	t.Data = append(t.Data, data)
}

func (t *MenuData) Get() []map[string]any {
	// Reset the menu
	var menu []map[string]any
	// Merge and append menus
	for _, appData := range t.Data {
		if t.PushData[appData.App] != nil {
			for _, datum := range t.PushData[appData.App].Data {
				appData.Data = append(appData.Data, datum)
			}
		}
	}

	sort.Slice(menu, func(i, j int) bool {
		return menu[i]["sort"].(int) < menu[j]["sort"].(int)
	})

	menu = getLoop(t.Prefix, "", t.Data)

	return menu
}

func getLoop(prefix, keys string, datas []*MenuData) []map[string]any {

	var list []map[string]any
	for _, items := range datas {
		key := strings.Join([]string{keys, items.Name}, "/")

		data := map[string]any{
			"name":     items.Name,
			"key":      key,
			"label":    items.Label,
			"route":    lo.Ternary[string](items.Route != "", prefix+"/"+items.Route, items.Route),
			"icon":     items.Icon,
			"meta":     items.Meta,
			"sort":     items.Sort,
			"children": []map[string]any{},
		}
		if len(items.Data) > 0 {
			data["children"] = getLoop(prefix, key, items.Data)
		}
		list = append(list, data)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["sort"].(int) < list[j]["sort"].(int)
	})

	return list
}
