package utils

import (
	"reflect"
	"strings"
)

// GetFieldNamesRecursive 递归获取f的成员名称。名称格式是全路径格式。
func GetFieldNamesRecursive(f reflect.StructField) []string {
	if f.Type.Kind() == reflect.Struct {
		names := make([]string, 0)

		for _, subfield := range GetExportedFields(f.Type) {
			for _, subName := range GetFieldNamesRecursive(subfield) {
				var fullNameBuilder strings.Builder
				fullNameBuilder.WriteString(f.Name)
				fullNameBuilder.WriteRune(rune('.'))
				fullNameBuilder.WriteString(subName)
				names = append(names, fullNameBuilder.String())
			}
		}
		return names
	} else {
		return []string{f.Name}
	}
}

// GetExportedFields 获取v中所有导出的成员
func GetExportedFields(v reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsExported() {
			fields = append(fields, field)
		}
	}
	return fields
}

// GetValueByPath 通过全路径获取结构体中这个字段的值，全路径以点号分割
func GetValueByPath(v reflect.Value, fullPath string) (interface{}, reflect.Kind) {
	for _, field := range strings.Split(fullPath, ".") {
		v = v.FieldByName(field)
	}
	return v.Interface(), v.Kind()
}
