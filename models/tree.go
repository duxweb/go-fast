package models

import (
	"gorm.io/gorm"
	"reflect"
)

type Tree interface {
	GetChildren() []*Tree
	GetId() uint
}

func ChildrenPreload(d *gorm.DB) *gorm.DB {
	return d.Preload("Children", ChildrenPreload)
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
