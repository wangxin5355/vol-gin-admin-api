package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/attribute_manager"
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

	//DelOnExecuting 删除前事件
	DelOnExecuting func([]any) *response.WebResponseContent
	//DelOnExecuted 删除后事件
	DelOnExecuted func([]any) *response.WebResponseContent
}

// 构造函数
func InitBaseService[T, T2 any](dbName string) *BaseService[T, T2] {
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

func (s *BaseService[T, T2]) GetDetailPage(options request.PageDataOptions) *response.PageGridData[map[string]any] {
	var master T2
	meta := attribute_manager.GetEntityMeta(master)
	if utils.IsNull(meta.DetailTableStr) {
		return &response.PageGridData[map[string]any]{Rows: nil, Total: 0}
	}
	return GetDetailPage(s.DB, meta, options)
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
func (s *BaseService[T, T2]) Del(c *gin.Context, keys []any) *response.WebResponseContent {
	return del[T, T2](c, s.DB, keys, s.DelOnExecuting, s.DelOnExecuted)
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

// GetDetailPage 获取明细表分页数据（按泛型 T 执行查询，不依赖外部 DetailTable 元数据）
func GetDetailPage(db *gorm.DB, detailEntityType attribute_manager.EntityMeta, options request.PageDataOptions) *response.PageGridData[map[string]any] {
	var total int64
	// 默认分页修正
	if options.Page <= 0 {
		options.Page = 1
	}
	if options.Rows <= 0 {
		options.Rows = 10
	}
	q := db.Table(detailEntityType.DetailTableStr)

	q = ApplyJsonWhereToDB(q, options)
	q = ApplyJsonSortToDB(q, options)

	//条件必须加上主表的主键条件
	q = q.Where(fmt.Sprintf("%s = ?", detailEntityType.Key), options.Value)

	countQuery := q.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return &response.PageGridData[map[string]any]{Rows: nil, Total: 0}
	}
	q = ApplyJsonPageToDB(q, options)
	var rows []map[string]any
	if err := q.Find(&rows).Error; err != nil {
		return &response.PageGridData[map[string]any]{Rows: nil, Total: 0}
	}
	return &response.PageGridData[map[string]any]{
		Total:   int(total),
		Rows:    rows,
		Summary: nil,
	}
}

// add 添加数据
func add[T, T2 any](c *gin.Context,
	db *gorm.DB,
	options request.SaveModel,
	AddOnExecuting, AddOnExecuted func(*T2) *response.WebResponseContent) *response.WebResponseContent {

	//entity = utils.DicToEntity[T2](options.MainData)
	entity := utils.MapToEntity[T2](options.MainData)
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
	//明细表处理
	//先查一下有没有
	detailData := options.DetailData
	if detailData != nil && len(detailData) > 0 {
		return addDetail[T2](c, db, &entity, options)
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
	res := map[string]any{
		"data": entity,
		"list": nil,
	}
	return response.Ok("添加成功", res)
}

// 添加明细
func addDetail[T2 any](c *gin.Context,
	db *gorm.DB,
	entity *T2, options request.SaveModel) *response.WebResponseContent {
	//获取实体信息
	//根据实体信息获取明细的实体
	var master T2
	meta := attribute_manager.GetEntityMeta(master)
	if utils.IsNull(meta.DetailTableStr) {
		return response.Ok("添加成功", entity)
	}

	tableName := meta.DetailTableStr
	var detailData, _ = utils.GetEntityListByTableName(tableName, options.DetailData)

	var userInfo = GetUserInfo(c)
	utils.SetDefaultValue[T2](entity, true, userInfo.UserID, userInfo.Username)
	//开启事务保存主表和明细表
	err := db.Transaction(func(tx *gorm.DB) error {
		//保存主表,并且拿到主键值赋给明细表
		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		pkName, pkValue, err := utils.GetPrimaryKey(tx, entity)
		if err != nil {
			return err
		}

		//由于是新增操作 所以不存在删除和修改 只需要处理新增
		for _, detailItem := range detailData {
			// 给明细表的主键字段赋值
			reflectVal := reflect.ValueOf(detailItem)

			// 如果是指针类型，获取其元素
			if reflectVal.Kind() == reflect.Ptr {
				reflectVal = reflectVal.Elem()
			}
			fieldVal := reflectVal.FieldByName(pkName)
			if !fieldVal.IsValid() {
				return fmt.Errorf("未找到字段: %s", pkName)
			}
			if !fieldVal.CanSet() {
				return fmt.Errorf("字段不可设置: %s", pkName)
			}
			if fieldVal.Kind() == reflect.Ptr {
				v := reflect.ValueOf(pkValue)
				ptr := reflect.New(v.Type())
				ptr.Elem().Set(v)
				fieldVal.Set(ptr)
			} else {
				fieldVal.Set(reflect.ValueOf(pkValue))
			}
			utils.SetDetailDefaultValue(detailItem, true, userInfo.UserID, userInfo.Username) // 默认值设置
			if err := tx.Create(detailItem).Error; err != nil {
				return fmt.Errorf("添加明细数据失败: %v", err)
			}
		}
		//提交事务
		return nil
	})
	if err != nil {
		return response.Error("添加失败: " + err.Error())
	}
	res := map[string]any{
		"data": entity,
		"list": detailData,
	}
	return response.Ok("添加成功", res)
}

// update 更新数据，只更新实体中存在的字段且排除主键
func update[T, T2 any](c *gin.Context,
	db *gorm.DB,
	options request.SaveModel,
	UpdateOnExecuting,
	UpdateOnExecuted func(*T2) *response.WebResponseContent) *response.WebResponseContent {

	// 用 DicToEntity[T2] 生成业务实体
	//entity := utils.DicToEntity[T2](options.MainData)
	entity := utils.MapToEntity[T2](options.MainData)

	//判断有明细的话直接走明细更新
	detailData := options.DetailData
	if detailData != nil && len(detailData) > 0 {
		return updateDetail[T2](c, db, &entity, options)
	}
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
	res := map[string]any{
		"data": entity,
		"list": nil,
	}
	return response.Ok("更新成功", res)
}

// update 明细
func updateDetail[T2 any](c *gin.Context,
	db *gorm.DB,
	entity *T2, options request.SaveModel) *response.WebResponseContent {
	//获取实体信息
	//根据实体信息获取明细的实体
	var master T2
	meta := attribute_manager.GetEntityMeta(master)
	if utils.IsNull(meta.DetailTableStr) {
		return response.Ok("更新成功", entity)
	}

	tableName := meta.DetailTableStr
	//把optnios.DetailData的数据转换为实例类
	var detailData, _ = utils.GetEntityListByTableName(tableName, options.DetailData)

	var userInfo = GetUserInfo(c)
	utils.SetDefaultValue[T2](entity, false, userInfo.UserID, userInfo.Username)
	//开启事务保存主表和明细表
	err := db.Transaction(func(tx *gorm.DB) error {
		//保存主表,并且拿到主键值赋给明细表
		if err := tx.Save(entity).Error; err != nil {
			return err
		}

		pkName, pkValue, err := utils.GetPrimaryKey(tx, entity)
		if err != nil {
			return err
		}

		//删除
		if len(options.DelKeys) > 0 {
			//根据明细表的主键key直接删除
			if err := db.Table(tableName).Delete(options.DelKeys).Error; err != nil {
				return fmt.Errorf("删除明细数据失败: " + err.Error())
			}
		}

		// 2. 处理新增和修改
		for _, detailItem := range detailData {
			// 判断该记录是新增还是修改
			reflectVal := reflect.ValueOf(detailItem)

			// 如果是指针类型，获取其元素
			if reflectVal.Kind() == reflect.Ptr {
				reflectVal = reflectVal.Elem()
			}
			fieldVal := reflectVal.FieldByName(pkName)
			if !fieldVal.IsValid() {
				return fmt.Errorf("未找到字段: %s", pkName)
			}
			if !fieldVal.CanSet() {
				return fmt.Errorf("字段不可设置: %s", pkName)
			}
			if fieldVal.Kind() == reflect.Ptr {
				v := reflect.ValueOf(pkValue)
				ptr := reflect.New(v.Type())
				ptr.Elem().Set(v)
				fieldVal.Set(ptr)
			} else {
				fieldVal.Set(reflect.ValueOf(pkValue))
			}
			utils.SetDetailDefaultValue(detailItem, false, userInfo.UserID, userInfo.Username) // 默认值设置
			if err := tx.Save(detailItem).Error; err != nil {
				return fmt.Errorf("修改明细数据失败: %v", err)
			}
		}
		//提交事务
		return nil
	})
	if err != nil {
		return response.Error("更新失败: " + err.Error())
	}
	//返回data中是  data=主表，list=明细表数据
	res := map[string]any{
		"data": entity,
		"list": detailData,
	}
	return response.Ok("更新成功", res)
}

// del 删除数据
func del[T, T2 any](c *gin.Context,
	db *gorm.DB,
	keys []any,
	DelOnExecuting, DelOnExecuted func([]any) *response.WebResponseContent) *response.WebResponseContent {
	if len(keys) == 0 {
		return response.Error("删除失败: 参数 keys 不能为空")
	}

	var entity T2
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
	// 删除前事件
	if DelOnExecuting != nil {
		beforeResp := DelOnExecuting(keys)
		if beforeResp.Status == false {
			return beforeResp
		}
	}
	// 删除后事件结果
	var afterResp *response.WebResponseContent
	//开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 执行删除
		if err := tx.Where(primaryField.DBName+" IN ?", keys).Delete(new(T2)).Error; err != nil {
			return err
		}
		//删除后事件
		if DelOnExecuted != nil {
			afterResp = DelOnExecuted(keys)
			if afterResp.Status == false {
				return fmt.Errorf(afterResp.Message)
			}
		}
		//提交事务
		return nil
	})
	if err != nil {
		return response.Error("删除失败: " + err.Error())
	}
	return response.Ok("删除成功", nil)
}

//-------------------------------
// 工具方法
//-------------------------------

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	data := utils.GetUserInfo(c)
	if data == nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在token")
		return nil
	}
	return data
}

// BindJsonToPageDataOptions 绑定分页参数
func BindJsonToPageDataOptions(c *gin.Context) (request.PageDataOptions, error) {
	var param request.PageDataOptions
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.WebResponse(response.Error(err.Error()), c)
		return param, err
	}
	return param, nil
}

// BindJsonToSaveModel 绑定保存参数
func BindJsonToSaveModel(c *gin.Context) (request.SaveModel, error) {
	var param request.SaveModel
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.WebResponse(response.Error(err.Error()), c)
		return param, err
	}
	return param, nil
}

// ShouldBindJSON 绑定JSON参数
func ShouldBindJSON[T any](c *gin.Context) (T, error) {
	var param T
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.WebResponse(response.Error(err.Error()), c)
		return param, err
	}
	return param, nil
}
