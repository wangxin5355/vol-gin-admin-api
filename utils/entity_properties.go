package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// DicToEntity 反射赋值
func DicToEntity[T any](data map[string]any) T {
	var entity T
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(&entity).Elem()

	for k, val := range data {
		assignFieldValue(t, v, k, val)
	}
	return entity
}

// 递归赋值
func assignFieldValue(t reflect.Type, v reflect.Value, k string, val any) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 递归处理嵌套结构体（匿名嵌套）
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if fieldValue.CanAddr() {
				assignFieldValue(field.Type, fieldValue.Addr().Elem(), k, val)
			} else {
				assignFieldValue(field.Type, fieldValue, k, val)
			}
			continue
		}

		if strings.EqualFold(field.Name, k) {
			if !fieldValue.CanSet() {
				continue
			}
			if field.Type.Kind() == reflect.Struct {
				if subMap, ok := val.(map[string]any); ok {
					subEntity := DicToEntityByType(subMap, field.Type)
					fieldValue.Set(reflect.ValueOf(subEntity).Elem()) // 赋值时取 Elem()
				}
				continue
			}
			converted, ok := ConvertType(val, field.Type)
			if ok {
				fieldValue.Set(reflect.ValueOf(converted))
			}
		}
	}
}

// 递归处理非泛型版本
func DicToEntityByType(data map[string]any, typ reflect.Type) any {
	entityPtr := reflect.New(typ) // 创建指针
	entity := entityPtr.Elem()

	for k, val := range data {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if strings.EqualFold(field.Name, k) {
				fieldValue := entity.Field(i)

				if !fieldValue.CanSet() {
					continue
				}

				if field.Type.Kind() == reflect.Struct {
					if subMap, ok := val.(map[string]any); ok {
						subEntity := DicToEntityByType(subMap, field.Type)
						fieldValue.Set(reflect.ValueOf(subEntity).Elem()) // 赋值时取 Elem()
					}
					continue
				}

				converted, ok := ConvertType(val, field.Type)
				if ok {
					fieldValue.Set(reflect.ValueOf(converted))
				}
			}
		}
	}
	return entityPtr.Interface() // 返回指针
}

// ConvertType 支持常见类型转换，保证赋值安全
func ConvertType(val any, typ reflect.Type) (any, bool) {
	if val == nil {
		return reflect.Zero(typ).Interface(), true
	}
	v := reflect.ValueOf(val)
	switch typ.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", val), true
	case reflect.Int:
		// 支持 int64/float64/string 转 int
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return int(v.Int()), true
		}
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return int(v.Float()), true
		}
		if v.Kind() == reflect.String {
			var i int
			_, err := fmt.Sscan(v.String(), &i)
			if err == nil {
				return i, true
			}
		}
	case reflect.Int8:
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return int8(v.Int()), true
		}
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return int8(v.Float()), true
		}
		if v.Kind() == reflect.String {
			var i int8
			_, err := fmt.Sscan(v.String(), &i)
			if err == nil {
				return i, true
			}
		}
	case reflect.Int16:
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return int16(v.Int()), true
		}
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return int16(v.Float()), true
		}
		if v.Kind() == reflect.String {
			var i int16
			_, err := fmt.Sscan(v.String(), &i)
			if err == nil {
				return i, true
			}
		}
	case reflect.Int32:
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return int32(v.Int()), true
		}
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return int32(v.Float()), true
		}
		if v.Kind() == reflect.String {
			var i int32
			_, err := fmt.Sscan(v.String(), &i)
			if err == nil {
				return i, true
			}
		}
	case reflect.Int64:
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return v.Int(), true
		}
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return int64(v.Float()), true
		}
		if v.Kind() == reflect.String {
			var i int64
			_, err := fmt.Sscan(v.String(), &i)
			if err == nil {
				return i, true
			}
		}
	case reflect.Float32:
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return float32(v.Float()), true
		}
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return float32(v.Int()), true
		}
		if v.Kind() == reflect.String {
			var f float32
			_, err := fmt.Sscan(v.String(), &f)
			if err == nil {
				return f, true
			}
		}
	case reflect.Float64:
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			return v.Float(), true
		}
		if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			return float64(v.Int()), true
		}
		if v.Kind() == reflect.String {
			var f float64
			_, err := fmt.Sscan(v.String(), &f)
			if err == nil {
				return f, true
			}
		}
	case reflect.Bool:
		if v.Kind() == reflect.Bool {
			return v.Bool(), true
		}
		if v.Kind() == reflect.String {
			if v.String() == "true" {
				return true, true
			} else if v.String() == "false" {
				return false, true
			}
		}
	}
	// 其他类型直接尝试转换
	if v.Type().ConvertibleTo(typ) {
		return v.Convert(typ).Interface(), true
	}
	return reflect.Zero(typ).Interface(), false
}

// SetDefaultValue 设置默认字段的值"CreateID", "Creator", "CreateDate"，"ModifyID", "Modifier", "ModifyDate"
func SetDefaultValue[T any](entity *T, isAdd bool, userID uint32, userName string) T {
	v := reflect.ValueOf(entity).Elem()
	setDefaultValueByReflect(v, isAdd, userID, userName)
	return *entity
}

// 非泛型递归处理嵌套结构体默认值
func setDefaultValueByReflect(v reflect.Value, isAdd bool, userID uint32, userName string) {
	t := v.Type()
	now := reflect.ValueOf(time.Now())
	uid := reflect.ValueOf(userID)
	uname := reflect.ValueOf(userName)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		// 递归处理嵌套结构体（包括匿名嵌套）
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Now()) {
			setDefaultValueByReflect(fieldValue, isAdd, userID, userName)
			continue
		}
		switch field.Name {
		case "CreateID", "Creator":
			if isAdd && (fieldValue.IsZero() || fieldValue.Interface() == uint32(0) || fieldValue.Interface() == "") {
				switch field.Type.Kind() {
				case reflect.Uint32, reflect.Uint, reflect.Uint64:
					fieldValue.Set(uid)
				case reflect.Int:
					fieldValue.Set(reflect.ValueOf(int(userID)))
				case reflect.Int32:
					fieldValue.Set(reflect.ValueOf(int32(userID)))
				case reflect.Int64:
					fieldValue.Set(reflect.ValueOf(int64(userID)))
				case reflect.String:
					fieldValue.Set(uname)
				}
			}
		case "CreateDate":
			if isAdd && (fieldValue.IsZero() || fieldValue.Interface() == "" || fieldValue.Interface() == nil) {
				if field.Type == reflect.TypeOf(time.Now()) {
					fieldValue.Set(now)
				}
			}
		case "ModifyID", "Modifier":
			if !isAdd && (fieldValue.IsZero() || fieldValue.Interface() == uint32(0) || fieldValue.Interface() == "") {
				switch field.Type.Kind() {
				case reflect.Uint32, reflect.Uint, reflect.Uint64:
					fieldValue.Set(uid)
				case reflect.Int:
					fieldValue.Set(reflect.ValueOf(int(userID)))
				case reflect.Int32:
					fieldValue.Set(reflect.ValueOf(int32(userID)))
				case reflect.Int64:
					fieldValue.Set(reflect.ValueOf(int64(userID)))
				case reflect.String:
					fieldValue.Set(uname)
				}
			}
		case "ModifyDate":
			if !isAdd && (fieldValue.IsZero() || fieldValue.Interface() == "" || fieldValue.Interface() == nil) {
				if field.Type == reflect.TypeOf(time.Now()) {
					fieldValue.Set(now)
				}
			}
		}
	}
}
