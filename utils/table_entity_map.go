package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
)

var tableEntityMap = map[string]any{
	"Sys_Dictionary":     &system.SysDictionary{},
	"Sys_DictionaryList": &system.SysDictionaryList{},
}

// GetEntityByTableName 根据表名获取实体类型
func GetEntityByTableName(tableName string) (any, error) {
	entity, ok := tableEntityMap[tableName]
	if !ok {
		return nil, fmt.Errorf("未找到表对应实体: %s", tableName)
	}
	t := reflect.TypeOf(entity).Elem()
	v := reflect.New(t)
	return v.Interface(), nil
}

// GetEntityListByTableName 根据表名获取实体类型并将数据转换为对应的实体list
func GetEntityListByTableName(tableName string, data []map[string]any) ([]any, error) {
	entity, err := GetEntityByTableName(tableName)
	if err != nil {
		return nil, err
	}

	elemType := reflect.TypeOf(entity).Elem()
	result := make([]any, 0, len(data))
	for _, m := range data {
		itemPtr := reflect.New(elemType).Interface()

		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName:              "json",
			Result:               itemPtr,
			WeaklyTypedInput:     true, // 自动转换
			IgnoreUntaggedFields: true, // 未定义字段忽略
			ZeroFields:           false,
			DecodeHook:           timeHook,
		})

		if err := decoder.Decode(m); err != nil {
			fmt.Printf("mapstructure 解析失败: %v\n数据: %+v\n", err, m)
			return nil, err
		}
		//	if err := mapToStruct(m, itemPtr); err != nil {
		//		fmt.Printf("map 转实体失败: %v, 原始数据=%v\n", err, m)
		//		return nil, err
		//	}

		result = append(result, itemPtr)
	}

	return result, nil
}

func timeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	// 目标是 *time.Time
	if to == reflect.TypeOf(&time.Time{}) {
		if data == nil || data == "" {
			return nil, nil // 不赋值
		}

		// mapstructure 默认传递 string 进来
		str, ok := data.(string)
		if !ok {
			return nil, nil // 不是 string 就跳过
		}

		// 尝试解析 RFC3339
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return nil, nil // 解析不了就跳过
		}

		return &t, nil
	}

	return data, nil
}

func mapToStruct(m map[string]any, ptr any) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("ptr 必须是指针类型")
	}

	v = v.Elem() // struct
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		jsonKey := field.Tag.Get("json")
		if jsonKey == "" {
			jsonKey = field.Name
		}

		if val, ok := m[jsonKey]; ok {
			err := setValue(fieldValue, val)
			if err != nil {
				return fmt.Errorf("字段 %s 赋值失败: %v (value=%v)", field.Name, err, val)
			}
		}
	}
	return nil
}
func setValue(field reflect.Value, val any) error {
	if !field.CanSet() {
		return fmt.Errorf("字段无法设置")
	}

	// ====== 处理指针类型 ======
	if field.Kind() == reflect.Pointer {
		// 处理 nil 值
		if val == nil {
			field.Set(reflect.Zero(field.Type()))
			return nil
		}

		elemType := field.Type().Elem()        // 指针指向的类型
		newVal := reflect.New(elemType).Elem() // 实例化一个新的值

		// 递归赋值（让它处理 time.Time、int、string 等）
		if err := setValue(newVal, val); err != nil {
			return err
		}

		field.Set(newVal.Addr())
		return nil
	}

	// ====== 处理 time.Time ======
	if field.Type() == reflect.TypeOf(time.Time{}) {
		switch v := val.(type) {
		case string:
			// 解析 RFC3339，例如 "2023-11-07T13:24:42+08:00"
			t, err := time.Parse(time.RFC3339, v)
			if err == nil {
				field.Set(reflect.ValueOf(t))
				return nil
			}

			// ⬆️如果你有其他时间格式，可以继续加

			return fmt.Errorf("time.Time 解析失败: %v", err)
		case time.Time:
			field.Set(reflect.ValueOf(v))
			return nil
		}
		return fmt.Errorf("无法将 %T 转为 time.Time", val)
	}

	// ====== 基础类型处理 ======
	switch field.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprintf("%v", val))
		return nil

	case reflect.Int, reflect.Int32, reflect.Int64:
		switch v := val.(type) {
		case float64:
			field.SetInt(int64(v))
		case string:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(n)
		case int, int32, int64:
			field.SetInt(reflect.ValueOf(v).Int())
		default:
			return fmt.Errorf("无法转换 %T 到 int", val)
		}
		return nil

	case reflect.Float32, reflect.Float64:
		switch v := val.(type) {
		case float64:
			field.SetFloat(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			field.SetFloat(f)
		default:
			return fmt.Errorf("无法转换 %T 到 float", val)
		}
		return nil

	case reflect.Bool:
		switch v := val.(type) {
		case bool:
			field.SetBool(v)
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			field.SetBool(b)
		default:
			return fmt.Errorf("无法转换 %T 到 bool", val)
		}
		return nil
	}

	return fmt.Errorf("不支持的字段类型: %s", field.Kind())
}
