package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/wangxin5355/vol-gin-admin-api/global"
	"gorm.io/gorm"
)

// 响应内容
type WebResponseContent struct {
	Status  bool
	Code    string
	Message string
	Data    any
}

func Ok(msg string, data any) *WebResponseContent {
	return &WebResponseContent{Status: true, Message: msg, Data: data}
}
func Error(msg string) *WebResponseContent {
	return &WebResponseContent{Status: false, Message: msg}
}

// 分页参数
type PageDataOptions struct {
	Sort   string             `json:"sort"`
	Order  string             `json:"order"`
	Page   int                `json:"page"`
	Rows   int                `json:"rows"`
	Export bool               `json:"-"`
	Filter []SearchParameters `json:"-"`
	Wheres string             `json:"wheres"`
	Value  any                `json:"value"`
}

// 分页结果
type PageGridData[T any] struct {
	Rows    []T
	Total   int
	Summary any
}

// SaveModel 示例
type SaveModel struct {
	MainData   map[string]any
	DetailData []map[string]any
	DelKeys    []any
}

// DicToEntity 反射赋值（简化版）
func DicToEntity[T any](data map[string]any) T {
	var entity T
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(&entity).Elem()
	for k, val := range data {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == k {
				v.Field(i).Set(reflect.ValueOf(val))
			}
		}
	}
	return entity
}

type SearchParameters struct {
	Name        string
	Value       string
	DisplayType string
}

// ApplyJsonWhereToDB 从参数转换为 GORM 查询条件
func ApplyJsonWhereToDB(db *gorm.DB, options PageDataOptions) *gorm.DB {
	jsonStr := options.Wheres
	var params []SearchParameters
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
func ApplyJsonSortToDB(db *gorm.DB, options PageDataOptions) *gorm.DB {
	if options.Sort == "" || options.Order == "" {
		return db
	}
	order := fmt.Sprintf("%s %s", options.Sort, options.Order)
	return db.Order(order)
}

// ApplyJsonPageToDB 分页语句解析
func ApplyJsonPageToDB(db *gorm.DB, options PageDataOptions) *gorm.DB {
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
func ApplyJsonToDB(db *gorm.DB, options PageDataOptions) *gorm.DB {
	db = ApplyJsonWhereToDB(db, options)
	db = ApplyJsonSortToDB(db, options)
	db = ApplyJsonPageToDB(db, options)
	return db
}

// GetPageData 传入一个实体，将其转换为 GORM 的映射对象
func GetPageData[T any](options PageDataOptions) *PageGridData[T] {
	var list []T
	var total int64
	// 获取 GORM DB 实例
	db := global.GVA_DB.Model(new(T))
	// 查询条件、排序、分页
	db = ApplyJsonToDB(db, options)

	// 先执行查询总数，如果是空的就不需要继续执行了
	if err := db.Count(&total).Error; err != nil {
		return &PageGridData[T]{Rows: nil, Total: 0}
	}
	// 执行查询
	if err := db.Find(&list).Error; err != nil {
		return &PageGridData[T]{Rows: nil, Total: 0}
	}
	return &PageGridData[T]{Rows: list, Total: int(total)}
}
