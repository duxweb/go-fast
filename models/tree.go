package models

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

func ChildrenPreload(d *gorm.DB) *gorm.DB {
	sort, ok := d.Get("tree_sort")
	q := d.Preload("Children", ChildrenPreload)
	if ok && sort != "" {
		q = q.Order(sort.(string) + " ASC")
	}
	return q
}

func ParentPreload(d *gorm.DB) *gorm.DB {
	return d.Preload("Parent", ParentPreload)
}

func CheckParentHas[T any](d *gorm.DB, id uint, parentID uint) bool {
	if id == parentID {
		return false
	}
	var models []T
	d.Preload("Children", ChildrenPreload).Where("parent_id = ?", id).Find(&models)
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
	query := `
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

	rows, err := db.Raw(fmt.Sprintf(query, table, table), id).Rows()
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
	query := `
	WITH RECURSIVE breadcrumbs AS (
		SELECT id, parent_id
		FROM %s
		WHERE id = ?

		UNION ALL

		SELECT t.id, t.parent_id
		FROM %s t
		INNER JOIN breadcrumbs b ON b.parent_id = t.id
	)
	SELECT id FROM breadcrumbs;`

	rows, err := db.Raw(fmt.Sprintf(query, table, table), id).Rows()
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
