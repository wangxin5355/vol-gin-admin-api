package utils

import (
	"fmt"
	"reflect"
)

// DicToEntity 反射赋值（简化版）
func DicToEntity[T any](data map[string]any) T {
	var entity T
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(&entity).Elem()
	for k, val := range data {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == k {
				fieldValue := v.Field(i)
				if fieldValue.CanSet() {
					converted, ok := ConvertType(val, field.Type)
					if ok {
						fieldValue.Set(reflect.ValueOf(converted))
					}
				}
			}
		}
	}
	return entity
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
