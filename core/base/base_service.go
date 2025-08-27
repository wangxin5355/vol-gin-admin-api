package base

import (
	"encoding/json"
	"fmt"
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
type BaseService[T, T2 any] struct {
	DB *gorm.DB
	// 查询前条件扩展
	QueryRelativeExpression func(*gorm.DB) *gorm.DB
	// 查询统计扩展
	SummaryExpress func(*gorm.DB) any
	// 查询后(从数据库查询的结果)
	GetPageDataOnExecuted func(*[]T)

	//AddOnExecuting 保存到数据库前事件
	AddOnExecuting func(*T2) *response.WebResponseContent
	//AddOnExecuted 保存到数据库后事件
	AddOnExecuted func(*T2) *response.WebResponseContent

	//编辑方法保存数据库前处理
	UpdateOnExecuting func(*T2) *response.WebResponseContent
	//编辑方法保存数据库后处理
	UpdateOnExecuted func(*T2) *response.WebResponseContent
}

// 构造函数
func NewBaseService[T, T2 any](dbName string) *BaseService[T, T2] {
	db := global.GetGlobalDBByDBName(dbName)
	if db == nil {
		panic("数据库连接未初始化或名称错误: " + dbName)
	}
	return &BaseService[T, T2]{
		DB: db,
	}
}

// getPageData 分页查询
func (s *BaseService[T, T2]) GetPageData(options request.PageDataOptions) *response.PageGridData[T] {
	return getPageData[T, T2](s.DB, options, s.QueryRelativeExpression, s.SummaryExpress, s.GetPageDataOnExecuted)
}

// add 添加
func (s *BaseService[T, T2]) Add(c *gin.Context, saveModel request.SaveModel) *response.WebResponseContent {
	return add[T, T2](c, s.DB, saveModel, s.AddOnExecuting, s.AddOnExecuted)
}

// update 更新
func (s *BaseService[T, T2]) Update(c *gin.Context, saveModel request.SaveModel) *response.WebResponseContent {
	return update[T, T2](c, s.DB, saveModel, s.UpdateOnExecuting, s.UpdateOnExecuted)
}

// del 删除
func (s *BaseService[T, T2]) Del(keys []any) *response.WebResponseContent {
	return del[T, T2](s.DB, keys)
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
func getPageData[T, T2 any](db *gorm.DB,
	options request.PageDataOptions,
	QueryRelativeExpression func(*gorm.DB) *gorm.DB,
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
	if QueryRelativeExpression != nil {
		db = QueryRelativeExpression(db)
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
func add[T, T2 any](c *gin.Context,
	db *gorm.DB,
	options request.SaveModel,
	AddOnExecuting, AddOnExecuted func(*T2) *response.WebResponseContent) *response.WebResponseContent {

	var entity T2
	entity = utils.DicToEntity[T2](options.MainData)
	var userInfo = GetUserInfo(c)
	if userInfo == nil {
		return response.Error("用户信息获取失败")
	}
	utils.SetDefaultValue[T2](&entity, true, userInfo.UserID, userInfo.Username)
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
func update[T, T2 any](c *gin.Context,
	db *gorm.DB,
	options request.SaveModel,
	UpdateOnExecuting,
	UpdateOnExecuted func(*T2) *response.WebResponseContent) *response.WebResponseContent {
	// 用 DicToEntity[T2] 生成业务实体
	entity := utils.DicToEntity[T2](options.MainData)
	var userInfo = GetUserInfo(c)
	utils.SetDefaultValue[T2](&entity, false, userInfo.UserID, userInfo.Username)

	// 解析结构体
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&entity); err != nil {
		return response.Error("更新失败: " + err.Error())
	}
	stmt.Dest = options.MainData

	// 获取主键字段及值
	primaryField := stmt.Schema.PrioritizedPrimaryField
	if primaryField == nil {
		return response.Error("更新失败: 未找到主键定义")
	}
	pkVal, hasPk := options.MainData[primaryField.Name]
	if !hasPk || pkVal == nil || pkVal == "" || pkVal == 0 {
		return response.Error("更新失败: 参数缺少主键字段或主键值为空")
	}

	if UpdateOnExecuting != nil {
		beforeResp := UpdateOnExecuting(&entity)
		if beforeResp.Status == false {
			return beforeResp
		}
	}

	// 构造更新 map：只包含有值的导出字段
	updateFields := utils.BuildEntityFields(entity, stmt)

	// 没有字段可更新
	if len(updateFields) == 0 {
		return response.Error("更新失败: 没有可更新的字段")
	}

	//保存后事件结果
	var afterResp *response.WebResponseContent
	// 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 执行更新
		if err := tx.Model(new(T)).
			Where(primaryField.DBName+" = ?", pkVal).
			Updates(updateFields).Error; err != nil {
			return err
		}
		//保存后事件
		if UpdateOnExecuted != nil {
			afterResp = UpdateOnExecuted(&entity)
			if afterResp.Status == false {
				return fmt.Errorf(afterResp.Message)
			}
		}
		//提交事务
		return nil
	})
	if err != nil {
		return response.Error("更新失败: " + err.Error())
	}
	return response.Ok("更新成功", entity)
}

// del 删除数据
func del[T, T2 any](db *gorm.DB, keys []any) *response.WebResponseContent {
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
