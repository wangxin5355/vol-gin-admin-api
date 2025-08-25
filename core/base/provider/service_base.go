package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	req "github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	res "github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"gorm.io/gorm"
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
					converted, ok := convertType(val, field.Type)
					if ok {
						fieldValue.Set(reflect.ValueOf(converted))
					}
				}
			}
		}
	}
	return entity
}

// convertType 支持常见类型转换，保证赋值安全
func convertType(val any, typ reflect.Type) (any, bool) {
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

// ApplyJsonWhereToDB 从参数转换为 GORM 查询条件
func ApplyJsonWhereToDB(db *gorm.DB, options req.PageDataOptions) *gorm.DB {
	jsonStr := options.Wheres
	var params []req.SearchParameters
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil || len(params) == 0 {
		return db
	}

	var whereParts []string
	var args []interface{}

	for _, p := range params {
		switch strings.ToLower(p.DisplayType) {
		case "equal":
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", p.Name))
			args = append(args, p.Value)
		case "like":
			whereParts = append(whereParts, fmt.Sprintf("%s LIKE ?", p.Name))
			args = append(args, "%"+p.Value+"%")
		case "greaterthan":
			whereParts = append(whereParts, fmt.Sprintf("%s > ?", p.Name))
			args = append(args, p.Value)
		case "lessthan":
			whereParts = append(whereParts, fmt.Sprintf("%s < ?", p.Name))
			args = append(args, p.Value)
		default:
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", p.Name))
			args = append(args, p.Value)
		}
	}

	where := strings.Join(whereParts, " AND ")
	if where == "" {
		return db
	}
	return db.Where(where, args...)
}

// ApplyJsonSortToDB 解析为排序语句
func ApplyJsonSortToDB(db *gorm.DB, options req.PageDataOptions) *gorm.DB {
	if options.Sort == "" || options.Order == "" {
		return db
	}
	order := fmt.Sprintf("%s %s", options.Sort, options.Order)
	return db.Order(order)
}

// ApplyJsonPageToDB 分页语句解析
func ApplyJsonPageToDB(db *gorm.DB, options req.PageDataOptions) *gorm.DB {
	if options.Page <= 0 {
		options.Page = 1
	}
	if options.Rows <= 0 {
		options.Rows = 10
	}
	offset := (options.Page - 1) * options.Rows
	return db.Offset(offset).Limit(options.Rows)
}

// ApplyJsonToDB 将参数转换为条件、排序、分页等数据
func ApplyJsonToDB(db *gorm.DB, options req.PageDataOptions) *gorm.DB {
	db = ApplyJsonWhereToDB(db, options)
	db = ApplyJsonSortToDB(db, options)
	db = ApplyJsonPageToDB(db, options)
	return db
}

// GetPageData 传入一个实体，将其转换为 GORM 的映射对象
func GetPageData[T any](db *gorm.DB, options req.PageDataOptions) *res.PageGridData[T] {
	var list []T
	var total int64
	// 获取 GORM DB 实例
	db = db.Model(new(T))
	// 查询条件、排序、分页
	db = ApplyJsonToDB(db, options)

	// 先执行查询总数，如果是空的就不需要继续执行了
	if err := db.Count(&total).Error; err != nil {
		return &res.PageGridData[T]{Rows: nil, Total: 0}
	}
	// 执行查询
	if err := db.Find(&list).Error; err != nil {
		return &res.PageGridData[T]{Rows: nil, Total: 0}
	}
	return &res.PageGridData[T]{Rows: list, Total: int(total)}
}

// Add 添加数据
func Add[T any](db *gorm.DB, options req.SaveModel) *res.WebResponseContent {
	var entity T
	// 反射赋值
	entity = DicToEntity[T](options.MainData)
	if err := db.Create(&entity).Error; err != nil {
		return res.Error("添加失败: " + err.Error())
	}
	return res.Ok("添加成功", entity)
}

// Update 更新数据，只更新实体中存在的字段且排除主键
func Update[T any](db *gorm.DB, options req.SaveModel) *res.WebResponseContent {
	var entity T
	entity = DicToEntity[T](options.MainData)

	// 解析结构体
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&entity); err != nil {
		return res.Error("更新失败: " + err.Error())
	}

	// 获取主键字段及值
	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return res.Error("更新失败: 未找到主键定义")
	}
	pkVal, hasPk := options.MainData[primaryField.Name]
	if !hasPk || pkVal == nil || pkVal == "" || pkVal == 0 {
		return res.Error("更新失败: 参数缺少主键字段或主键值为空")
	}

	// 构造更新 map：只包含实体字段且非主键
	t := reflect.TypeOf(entity)
	updateFields := make(map[string]any)
	for k, v := range options.MainData {
		found := false
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// 字段名匹配
			if field.Name == k {
				found = true
				break
			}
		}
		if !found {
			continue // 跳过不在实体中的字段
		}

		// 跳过主键字段
		skip := false
		for _, pf := range stmt.Schema.PrimaryFields {
			if k == pf.Name || k == pf.DBName ||
				strings.EqualFold(k, pf.Name) || strings.EqualFold(k, pf.DBName) {
				skip = true
				break
			}
		}
		if !skip {
			updateFields[k] = v
		}
	}

	// 没有字段可更新
	if len(updateFields) == 0 {
		return res.Error("更新失败: 没有可更新的字段")
	}

	// 更新数据库
	if err := db.Model(new(T)).
		Where(primaryField.DBName+" = ?", pkVal).
		Updates(updateFields).Error; err != nil {
		return res.Error("更新失败: " + err.Error())
	}

	return res.Ok("更新成功", entity)
}

// Del 删除数据
func Del[T any](db *gorm.DB, keys []any) *res.WebResponseContent {
	if len(keys) == 0 {
		return res.Error("删除失败: 参数 keys 不能为空")
	}

	var entity T
	// 解析结构体
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&entity); err != nil {
		return res.Error("删除失败: " + err.Error())
	}

	// 获取主键字段
	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return res.Error("删除失败: 未找到主键定义")
	}

	// 执行删除
	if err := db.Where(primaryField.DBName+" IN ?", keys).Delete(new(T)).Error; err != nil {
		return res.Error("删除失败: " + err.Error())
	}

	return res.Ok("删除成功", nil)
}
