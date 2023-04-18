package helper

import (
	"github.com/demdxx/gocast/v2"
)

func SliceToTree(data []map[string]any, idField string, pidField string, sonField string) []map[string]any {
	var dataMap = make(map[uint]map[string]any)
	tree := []map[string]any{}
	for i, datum := range data {
		var id = gocast.Number[uint](datum[idField])
		dataMap[id] = datum
		data[i][sonField] = []map[string]any{}
	}
	for _, datum := range data {
		var pid = gocast.Number[uint](datum[pidField])
		if pid != 0 {
			dataMap[pid][sonField] = append(dataMap[pid][sonField].([]map[string]any), datum)
		} else {
			tree = append(tree, datum)
		}
	}
	return tree
}

func GetTreeNode(data any, id uint, idField string, sonField string) map[string]any {
	for _, datum := range data.([]map[string]any) {
		if gocast.Number[uint](datum[idField]) == id {
			return datum
		}
		if _, ok := datum[sonField]; ok {
			result := GetTreeNode(datum[sonField].([]map[string]any), id, idField, sonField)
			if len(result) > 0 {
				return result
			}
		}
	}
	return map[string]any{}
}
