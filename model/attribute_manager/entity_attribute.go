package attribute_manager

import (
	"reflect"
	"strings"
)

// EntityMeta 存储实体的元数据信息
type EntityMeta struct {
	TableCnName       string
	TableName         string
	DetailTable       []reflect.Type
	DetailTableStr    string
	DetailTableCnName string
	DBServer          string
	Key               string //主表和明细表的关联键
}

// GetEntityMeta 读取标签并解析到 EntityMeta 结构体中
func GetEntityMeta(v interface{}) EntityMeta {
	meta := EntityMeta{}

	t := reflect.TypeOf(v)
	if t == nil {
		return meta
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 只处理结构体且至少有一个字段（通常第一个字段为 `_ struct{} `entity:"..."`）
	if t.Kind() != reflect.Struct || t.NumField() == 0 {
		return meta
	}

	tag := t.Field(0).Tag.Get("entity")
	attrs := parseEntityTag(tag)

	// 根据常见键名填充 EntityMeta（不区分大小写）
	if v, ok := attrs["table_name"]; ok && v != "" {
		meta.TableName = v
	} else if v, ok := attrs["tablename"]; ok && v != "" {
		meta.TableName = v
	}

	if v, ok := attrs["table_cn_name"]; ok && v != "" {
		meta.TableCnName = v
	} else if v, ok := attrs["tablecnname"]; ok && v != "" {
		meta.TableCnName = v
	}

	if v, ok := attrs["detail_table_cn_name"]; ok && v != "" {
		meta.DetailTableCnName = v
	} else if v, ok := attrs["detailtablecnname"]; ok && v != "" {
		meta.DetailTableCnName = v
	}

	if v, ok := attrs["db_server"]; ok && v != "" {
		meta.DBServer = v
	} else if v, ok := attrs["dbserver"]; ok && v != "" {
		meta.DBServer = v
	}

	if v, ok := attrs["detail_table"]; ok && v != "" {
		meta.DetailTableStr = v
	} else if v, ok := attrs["detailtable"]; ok && v != "" {
		meta.DetailTableStr = v
	}
	if v, ok := attrs["key"]; ok && v != "" {
		meta.Key = v
	} else if v, ok := attrs["Key"]; ok && v != "" {
		meta.Key = v
	}

	return meta
}

// GetEntityAttribute 读取标签并解析为 map[string]string
func GetEntityAttribute(v interface{}) map[string]string {
	t := reflect.TypeOf(v)
	if t == nil {
		return map[string]string{}
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct || t.NumField() == 0 {
		return map[string]string{}
	}
	tag := t.Field(0).Tag.Get("entity")
	return parseEntityTag(tag)
}

// parseEntityTag 将 "k=v;k2=v2" 解析为 map，键统一为小写并去除空白
func parseEntityTag(tag string) map[string]string {
	result := map[string]string{}
	if tag == "" {
		return result
	}

	for _, part := range split(tag, ';') {
		kv := split(part, '=')
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		value := strings.TrimSpace(kv[1])
		if key != "" {
			result[key] = value
		}
	}
	return result
}

func split(s string, sep rune) []string {
	var res []string
	curr := ""
	for _, ch := range s {
		if ch == sep {
			res = append(res, curr)
			curr = ""
		} else {
			curr += string(ch)
		}
	}
	if curr != "" {
		res = append(res, curr)
	}
	return res
}
