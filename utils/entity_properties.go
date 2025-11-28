package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
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

// SetEntityDefaultValue 从 gin.Context 中获取用户信息并设置默认值
func SetEntityDefaultValue(v reflect.Value, c *gin.Context, isAdd bool) {
	data := GetUserInfo(c)
	if data == nil {
		return
	}
	userID := data.UserID
	userName := data.Username

	SetDetailDefaultValue(v.Interface(), isAdd, userID, userName)
}

// SetDetailDefaultValue 给明细表的集合附上默认值 接收detailItem.Addr().Interface()
func SetDetailDefaultValue(detailItem any, isAdd bool, userID uint32, userName string) {
	v := reflect.ValueOf(detailItem)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		setDefaultValueByReflect(v, isAdd, userID, userName)
	case reflect.Map:
		setDefaultValueForMap(v, isAdd, userID, userName)
	default:
		fmt.Println("不支持的类型:", v.Kind())
	}
}
func setDefaultValueForMap(v reflect.Value, isAdd bool, userID uint32, userName string) {
	if v.Type().Key().Kind() != reflect.String {
		fmt.Println("只支持 map[string]interface{} 类型")
		return
	}

	if isAdd {
		v.SetMapIndex(reflect.ValueOf("Creator"), reflect.ValueOf(userName))
		v.SetMapIndex(reflect.ValueOf("CreateDate"), reflect.ValueOf(time.Now()))
		v.SetMapIndex(reflect.ValueOf("CreateID"), reflect.ValueOf(userID))
	} else {
		v.SetMapIndex(reflect.ValueOf("Modifier"), reflect.ValueOf(userName))
		v.SetMapIndex(reflect.ValueOf("ModifyDate"), reflect.ValueOf(time.Now()))
		v.SetMapIndex(reflect.ValueOf("ModifyID"), reflect.ValueOf(userID))
	}
}

// 判断字段是否存在（支持嵌套结构体）
func hasFieldRecursive(t reflect.Type, name string) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 直接匹配字段名
		if strings.EqualFold(field.Name, name) {
			return true
		}

		// 如果是嵌套结构体，递归进去
		if field.Type.Kind() == reflect.Struct {
			if hasFieldRecursive(field.Type, name) {
				return true
			}
		}
	}
	return false
}

// 判断值是否为空
func IsEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	// 处理 time.Time 类型
	if t, ok := v.(time.Time); ok {
		return t.IsZero()
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	case reflect.Struct:
		//时间类型
		if rv.Type() == reflect.TypeOf(time.Time{}) {
			return rv.Interface().(time.Time).IsZero()
		}
	}
	return false
}

// BuildEntityFields 构造 updateFields，只取匿名结构体的字段
func BuildEntityFields(entity any, stmt *gorm.Statement) map[string]any {
	result := map[string]any{}

	// mapstructure 默认不会递归匿名结构体，需要设置 Squash 标签或者 DecodeHook
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json", // 使用 json tag 作为 key
		Result:           &result,
		WeaklyTypedInput: true, // 自动类型转换
	})

	_ = decoder.Decode(entity)
	return result

	//updateFields := make(map[string]any)
	//
	//v := reflect.ValueOf(entity)
	//if v.Kind() == reflect.Ptr {
	//	v = v.Elem()
	//}
	//t := v.Type()
	//
	//for i := 0; i < t.NumField(); i++ {
	//	field := t.Field(i)
	//	fieldName := field.Name
	//	// 跳过主键字段
	//	skip := false
	//	for _, pf := range stmt.Schema.PrimaryFields {
	//		if strings.EqualFold(fieldName, pf.Name) || strings.EqualFold(fieldName, pf.DBName) {
	//			skip = true
	//			break
	//		}
	//	}
	//	if skip {
	//		continue
	//	}
	//	// 递归处理匿名结构体
	//	if field.Anonymous && field.Type.Kind() == reflect.Struct {
	//		anonVal := v.Field(i)
	//		anonType := field.Type
	//		for j := 0; j < anonType.NumField(); j++ {
	//			anonField := anonType.Field(j)
	//			anonFieldName := anonField.Name
	//			skipAnon := false
	//			for _, pf := range stmt.Schema.PrimaryFields {
	//				if strings.EqualFold(anonFieldName, pf.Name) || strings.EqualFold(anonFieldName, pf.DBName) {
	//					skipAnon = true
	//					break
	//				}
	//			}
	//			if skipAnon {
	//				continue
	//			}
	//			anonFieldVal := anonVal.Field(j).Interface()
	//			if !IsEmptyValue(anonFieldVal) {
	//				updateFields[anonFieldName] = anonFieldVal
	//			}
	//		}
	//		continue
	//	}
	//	// 普通字段，只加入非零/非空值
	//	fieldVal := v.Field(i).Interface()
	//	if !IsEmptyValue(fieldVal) {
	//		updateFields[fieldName] = fieldVal
	//	}
	//}
	//return updateFields
}

// JsonToEntity json转换为实体
func JsonToEntity[T any](jsonStr string) T {
	var entity T
	json.Unmarshal([]byte(jsonStr), &entity)
	return entity
}

// MapToEntity 自动将 map[string]any 转为泛型实体 T（支持自动类型转换）
func MapToEntity[T any](data map[string]any) T {
	var entity T
	val := reflect.ValueOf(&entity).Elem() // struct
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// json tag 优先
		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		raw, ok := data[key]
		if !ok || raw == nil {
			continue
		}

		setFieldValue(fieldVal, raw)
	}

	return entity
}

// 将 map 中的 value 设置到字段中（自动类型转换）
func setFieldValue(field reflect.Value, raw any) {
	if !field.CanSet() {
		return
	}

	// ---- 指针类型处理 ----
	if field.Kind() == reflect.Pointer {
		if raw == nil {
			field.Set(reflect.Zero(field.Type()))
			return
		}
		elem := reflect.New(field.Type().Elem()).Elem()
		setFieldValue(elem, raw)
		field.Set(elem.Addr())
		return
	}

	// ---- time.Time 特殊处理 ----
	if field.Type() == reflect.TypeOf(time.Time{}) {
		switch v := raw.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, v)
			if err == nil {
				field.Set(reflect.ValueOf(t))
			}
		case time.Time:
			field.Set(reflect.ValueOf(v))
		}
		return
	}

	switch field.Kind() {

	case reflect.String:
		field.SetString(fmt.Sprintf("%v", raw))

	case reflect.Int, reflect.Int32, reflect.Int64:
		switch v := raw.(type) {
		case float64:
			field.SetInt(int64(v))
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				field.SetInt(i)
			}
		case int, int32, int64:
			field.SetInt(reflect.ValueOf(v).Int())
		}

	case reflect.Float32, reflect.Float64:
		switch v := raw.(type) {
		case float64:
			field.SetFloat(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				field.SetFloat(f)
			}
		}

	case reflect.Bool:
		switch v := raw.(type) {
		case bool:
			field.SetBool(v)
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				field.SetBool(b)
			}
		}
	}
}

// Contains 判断字符串数组是否包含某个字符串
func Contains(arr []any, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// GetPrimaryKey 获取实体的主键字段名及值
func GetPrimaryKey(db *gorm.DB, entity any) (string, any, error) {
	if db == nil {
		return "", nil, fmt.Errorf("db cannot be nil")
	}
	if entity == nil {
		return "", nil, fmt.Errorf("entity cannot be nil")
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(entity); err != nil {
		return "", nil, fmt.Errorf("parse entity failed: %w", err)
	}

	if stmt.Schema.PrioritizedPrimaryField == nil {
		return "", nil, fmt.Errorf("no primary key defined")
	}

	pf := stmt.Schema.PrioritizedPrimaryField

	v := reflect.ValueOf(entity)

	// 处理 interface 包裹
	for v.Kind() == reflect.Interface {
		if v.IsNil() {
			return "", nil, fmt.Errorf("entity is nil interface")
		}
		v = v.Elem()
	}

	// 处理指针
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "", nil, fmt.Errorf("entity is nil pointer")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("entity is not struct or pointer to struct")
	}

	// 尝试通过字段名获取
	field := v.FieldByName(pf.Name)
	if !field.IsValid() {
		// 尝试通过 DBName 获取
		if f, ok := stmt.Schema.FieldsByDBName[pf.DBName]; ok {
			field = v.FieldByName(f.Name)
		}
	}
	if !field.IsValid() {
		return "", nil, fmt.Errorf("primary key field %s not found", pf.Name)
	}

	var value any
	if field.Kind() == reflect.Pointer && field.IsNil() {
		value = nil
	} else {
		value = field.Interface()
	}

	return pf.DBName, value, nil
}
