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

func GetTreeNode[T comparable](data []map[string]any, id T, idField string, sonField string) map[string]any {
	for _, datum := range data {
		if datum[idField].(T) == id {
			return datum
		}
		if _, ok := datum[sonField]; ok {
			result := GetTreeNode[T](datum[sonField].([]map[string]any), id, idField, sonField)
			if len(result) > 0 {
				return result
			}
		}
	}
	return map[string]any{}
}

// 获取树形父节点
func GetTreeParentNode[T comparable](data []map[string]any, id T, idField string, sonField string) map[string]any {
	for _, datum := range data {
		if _, ok := datum[sonField]; ok {
			sons := datum[sonField].([]map[string]any)
			for _, son := range sons {
				if son[idField].(T) == id {
					return datum
				}
				result := GetTreeParentNode[T](sons, id, idField, sonField)
				if len(result) > 0 {
					return result
				}
			}
		}
	}
	return map[string]any{}

}
