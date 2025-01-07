package models

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/go-errors/errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type TreePreloadConfig struct {
	Sort     string
	Preloads []string
}

func TreePreload(config TreePreloadConfig) func(d *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		q := d.Preload("Children", TreePreload(config))

		// 处理其他预加载
		for _, preload := range config.Preloads {
			if preload != "Children" { // 避免重复预加载Children
				q = q.Preload(preload)
				fmt.Println(preload)
			}
		}

		// 处理排序
		if config.Sort != "" {
			q = q.Order(config.Sort + " ASC")
		}
		return q
	}
}

func ParentPreload(d *gorm.DB) *gorm.DB {
	return d.Preload("Parent", ParentPreload)
}

func CheckParentHas[T any](d *gorm.DB, id uint, parentID uint) bool {
	if id == parentID {
		return false
	}
	var models []T
	d.Preload("Children", TreePreload(TreePreloadConfig{
		Sort: "sort",
	})).Where("parent_id = ?", id).Find(&models)
	return checkParent(models, parentID)
}

func checkParent[T any](data []T, parentID uint) bool {
	for _, item := range data {
		id := getFieldValue(item, "ID")
		if id == nil {
			continue
		}
		if id.(uint) == parentID {
			return false
		}
		childrenValue := getFieldValue(item, "Children")
		if childrenValue == nil {
			continue
		}
		models := childrenValue.([]T)
		return checkParent[T](models, parentID)
	}
	return true
}

func getFieldValue(s any, fieldName string) any {
	sValue := reflect.ValueOf(s)
	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}
	if sValue.Kind() != reflect.Struct {
		return nil
	}
	fieldValue := sValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil
	}

	return fieldValue.Interface()
}

func QuerySubIDs(model any, id uint, db *gorm.DB) ([]uint, error) {
	var results []uint
	t := reflect.TypeOf(model)
	table := db.NamingStrategy.TableName(t.Name())

	var query string
	dialectName := db.Dialector.Name()

	switch dialectName {
	case "dm", "oracle":
		// 达梦和Oracle使用 CONNECT BY 语法
		query = `
		SELECT id
		FROM %s
		START WITH id = ?
		CONNECT BY PRIOR id = parent_id`
	default:
		// MySQL, PostgreSQL 等使用 WITH RECURSIVE 语法
		query = `
		WITH RECURSIVE subordinates AS (
			SELECT id
			FROM %s
			WHERE id = ?

			UNION ALL

			SELECT t.id
			FROM %s t
			INNER JOIN subordinates s ON t.parent_id = s.id
		)
		SELECT id FROM subordinates;`
	}

	var rows *sql.Rows
	var err error
	if dialectName == "dm" || dialectName == "oracle" {
		rows, err = db.Raw(fmt.Sprintf(query, table), id).Rows()
	} else {
		rows, err = db.Raw(fmt.Sprintf(query, table, table), id).Rows()
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemId uint
		if err := rows.Scan(&itemId); err != nil {
			return nil, err
		}
		results = append(results, itemId)
	}

	return results, nil
}

func QueryParentIDs(model any, id uint, db *gorm.DB) ([]uint, error) {
	var results []uint
	t := reflect.TypeOf(model)
	table := db.NamingStrategy.TableName(t.Name())

	var query string
	dialectName := db.Dialector.Name()

	switch dialectName {
	case "dm", "oracle":
		// 达梦和Oracle使用 CONNECT BY 语法
		query = `
		SELECT id
		FROM %s
		START WITH id = ?
		CONNECT BY PRIOR parent_id = id`
	default:
		// MySQL, PostgreSQL 等使用 WITH RECURSIVE 语法
		query = `WITH RECURSIVE breadcrumbs AS (
			SELECT id, parent_id
			FROM %s
			WHERE id = ?

			UNION ALL

			SELECT t.id, t.parent_id
			FROM %s t
			INNER JOIN breadcrumbs b ON b.parent_id = t.id
		)
		SELECT id FROM breadcrumbs`
	}

	var rows *sql.Rows
	var err error
	if dialectName == "dm" || dialectName == "oracle" {
		rows, err = db.Raw(fmt.Sprintf(query, table), id).Rows()
	} else {
		rows, err = db.Raw(fmt.Sprintf(query, table, table), id).Rows()
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemId uint
		if err := rows.Scan(&itemId); err != nil {
			return nil, err
		}
		results = append(results, itemId)
	}

	return results, nil
}

func TreeSort(db *gorm.DB, model any, id uint64, beforeId uint64, parentId uint64) error {

	// 获取当前节点
	currentNode := map[string]any{}
	if err := db.Model(model).First(&currentNode, id).Error; err != nil {
		return errors.New("node not found")
	}

	// 更新父节点ID
	nodeParentID := lo.Ternary(parentId > 0, &parentId, nil)

	err := db.Model(model).Where("id = ?", id).Update("parent_id", nodeParentID).Error
	if err != nil {
		return errors.New("update parent id failed")
	}

	// 查询所有子节点
	siblings := []map[string]any{}
	query := db.Model(model)
	if parentId > 0 {
		query = query.Where("parent_id = ?", parentId)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.Order("sort asc").Find(&siblings).Error; err != nil {
		return errors.New("get siblings failed")
	}

	// 更新子节点顺序
	newSiblings := []map[string]any{}
	for _, sibling := range siblings {
		if sibling["id"] == uint(beforeId) {
			newSiblings = append(newSiblings, sibling)
			newSiblings = append(newSiblings, currentNode)
		} else if sibling["id"] != uint(id) {
			newSiblings = append(newSiblings, sibling)
		}
	}

	// 如果beforeID为0,则插入到最前面
	if beforeId == 0 {
		newSiblings = append([]map[string]any{currentNode}, newSiblings...)
	}

	// 重新排序所有同级节点
	for i, sibling := range newSiblings {
		if err := db.Model(model).Where("id = ?", sibling["id"]).Update("sort", i).Error; err != nil {
			return errors.New("update child node sort failed")
		}
	}

	return nil
}
