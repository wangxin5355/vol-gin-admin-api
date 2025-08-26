package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	systemReq "github.com/wangxin5355/vol-gin-admin-api/model/system/request"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"gorm.io/gorm"
)

// BaseService 泛型基类
type BaseService[T any] struct {
	DB *gorm.DB
	// 查询前条件扩展
	QueryRelativeExpression func(*gorm.DB) *gorm.DB
	// 查询统计扩展
	SummaryExpress func(*gorm.DB) any
	// 查询后(从数据库查询的结果)
	GetPageDataOnExecuted func(*[]T)

	//AddOnExecuting 保存到数据库前事件
	AddOnExecuting func(*T) *response.WebResponseContent
	//AddOnExecuted 保存到数据库后事件
	AddOnExecuted func(*T) *response.WebResponseContent
}

// 构造函数
func NewBaseService[T any](dbName string) *BaseService[T] {
	db := global.GetGlobalDBByDBName(dbName)
	if db == nil {
		panic("数据库连接未初始化或名称错误: " + dbName)
	}
	return &BaseService[T]{
		DB: db,
	}
}

// getPageData 分页查询
func (s *BaseService[T]) GetPageData(options request.PageDataOptions) *response.PageGridData[T] {
	return getPageData[T](s.DB, options, s.QueryRelativeExpression, s.SummaryExpress, s.GetPageDataOnExecuted)
}

// add 添加
func (s *BaseService[T]) Add(c *gin.Context, saveModel request.SaveModel) *response.WebResponseContent {
	return add[T](c, s.DB, saveModel, s.AddOnExecuting, s.AddOnExecuted)
}

// update 更新
func (s *BaseService[T]) Update(saveModel request.SaveModel) *response.WebResponseContent {
	return update[T](s.DB, saveModel)
}

// del 删除
func (s *BaseService[T]) Del(keys []any) *response.WebResponseContent {
	return del[T](s.DB, keys)
}

//--------------------------------------------------------------
//具体实现
//--------------------------------------------------------------

// ApplyJsonWhereToDB 从参数转换为 GORM 查询条件
func ApplyJsonWhereToDB(db *gorm.DB, options request.PageDataOptions) *gorm.DB {
	jsonStr := options.Wheres
	var params []request.SearchParameters
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
func ApplyJsonSortToDB(db *gorm.DB, options request.PageDataOptions) *gorm.DB {
	if options.Sort == "" || options.Order == "" {
		return db
	}
	order := fmt.Sprintf("%s %s", options.Sort, options.Order)
	return db.Order(order)
}

// ApplyJsonPageToDB 分页语句解析
func ApplyJsonPageToDB(db *gorm.DB, options request.PageDataOptions) *gorm.DB {
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
func ApplyJsonToDB(db *gorm.DB, options request.PageDataOptions) *gorm.DB {
	db = ApplyJsonWhereToDB(db, options)
	db = ApplyJsonSortToDB(db, options)
	db = ApplyJsonPageToDB(db, options)
	return db
}

// getPageData 传入一个实体，将其转换为 GORM 的映射对象
func getPageData[T any](db *gorm.DB,
	options request.PageDataOptions,
	queryRelativeExpression func(*gorm.DB) *gorm.DB,
	SummaryExpress func(*gorm.DB) any,
	GetPageDataOnExecuted func(*[]T)) *response.PageGridData[T] {
	var list []T
	var total int64
	// 获取 GORM DB 实例
	db = db.Model(new(T))
	// 定义返回类
	var res = &response.PageGridData[T]{Rows: nil, Total: 0}
	// 查询条件、排序、分页
	db = ApplyJsonToDB(db, options)
	// 查询前条件扩展
	if queryRelativeExpression != nil {
		db = queryRelativeExpression(db)
	}
	// 先执行查询总数，如果是空的就不需要继续执行了
	if err := db.Count(&total).Error; err != nil {
		return res
	}
	// 执行查询
	if err := db.Find(&list).Error; err != nil {
		return res
	}
	// 统计扩展
	if SummaryExpress != nil {
		res.Summary = SummaryExpress(db)
	}
	// 查询后事件(从数据库查询的结果)
	if GetPageDataOnExecuted != nil {
		GetPageDataOnExecuted(&list)
	}
	res.Rows = list
	res.Total = int(total)
	return res
}

// add 添加数据
func add[T any](c *gin.Context,
	db *gorm.DB,
	options request.SaveModel,
	AddOnExecuting, AddOnExecuted func(*T) *response.WebResponseContent) *response.WebResponseContent {

	var entity T
	entity = utils.DicToEntity[T](options.MainData)
	var userInfo = GetUserInfo(c)
	utils.SetDefaultValue[T](&entity, true, userInfo.UserID, userInfo.Username)
	// 保存前事件
	if AddOnExecuting != nil {
		beforeResp := AddOnExecuting(&entity)
		if beforeResp.Status == false {
			return beforeResp
		}
	}
	// 保存后事件结果
	var afterResp *response.WebResponseContent
	// 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}
		//保存后事件
		if AddOnExecuted != nil {
			afterResp = AddOnExecuted(&entity)
			if afterResp.Status == false {
				return fmt.Errorf(afterResp.Message)
			}
		}
		//提交事务
		return nil
	})
	if err != nil {
		return response.Error("添加失败: " + err.Error())
	}
	return response.Ok("添加成功", entity)
}

// update 更新数据，只更新实体中存在的字段且排除主键
func update[T any](db *gorm.DB, options request.SaveModel) *response.WebResponseContent {
	var entity T
	entity = utils.DicToEntity[T](options.MainData)

	// 解析结构体
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&entity); err != nil {
		return response.Error("更新失败: " + err.Error())
	}

	// 获取主键字段及值
	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return response.Error("更新失败: 未找到主键定义")
	}
	pkVal, hasPk := options.MainData[primaryField.Name]
	if !hasPk || pkVal == nil || pkVal == "" || pkVal == 0 {
		return response.Error("更新失败: 参数缺少主键字段或主键值为空")
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
		return response.Error("更新失败: 没有可更新的字段")
	}

	// 更新数据库
	if err := db.Model(new(T)).
		Where(primaryField.DBName+" = ?", pkVal).
		Updates(updateFields).Error; err != nil {
		return response.Error("更新失败: " + err.Error())
	}

	return response.Ok("更新成功", entity)
}

// del 删除数据
func del[T any](db *gorm.DB, keys []any) *response.WebResponseContent {
	if len(keys) == 0 {
		return response.Error("删除失败: 参数 keys 不能为空")
	}

	var entity T
	// 解析结构体
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&entity); err != nil {
		return response.Error("删除失败: " + err.Error())
	}

	// 获取主键字段
	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return response.Error("删除失败: 未找到主键定义")
	}

	// 执行删除
	if err := db.Where(primaryField.DBName+" IN ?", keys).Delete(new(T)).Error; err != nil {
		return response.Error("删除失败: " + err.Error())
	}

	return response.Ok("删除成功", nil)
}

// 获取用户信息
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	data := utils.GetUserInfo(c)
	if data == nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在token")
		return nil
	}
	return data
}
