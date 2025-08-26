package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/core/base"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system/partial"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

//
//type TestService struct {
//}
//
//func Db() *gorm.DB {
//	return global.GetGlobalDBByDBName(string(initialize.DbGin))
//}
//
//// getPageData 获取分页数据
//func (s *TestService) getPageData(options request.PageDataOptions) *response.PageGridData[system.SysUser] {
//	// 调用父类方法：return s.ServiceBase.getPageData(options)
//	//limit := options.Rows
//	//offset := (options.Page - 1) * limit
//	//db := global.GVA_DB.Model(&system.SysUser{})
//
//	//var userLise []system.SysUser
//	//options.Wheres = "[{\"name\":\"Username\",\"value\":\"u\",\"displayType\":\"like\"}]"
//	//db = provider.ApplyJsonToDB(db, options)
//	////db.Limit(limit).Offset(offset).Find(&userLise)
//	//db.Find(&userLise)
//	//return &provider.PageGridData[system.SysUser]{
//	//	Rows:  userLise,
//	//	Total: 100,
//	return provider.getPageData[system.SysUser](Db(), options)
//}
//
//// add 添加方法
//func (s *TestService) add(saveModel request.SaveModel) *response.WebResponseContent {
//	return provider.add[system.SysUser](Db(), saveModel)
//}
//
//// update 更新方法
//func (s *TestService) update(saveModel request.SaveModel) *response.WebResponseContent {
//	return provider.update[system.SysUser](Db(), saveModel)
//}
//
//// del 删除方法
//func (s *TestService) del(keys []any) *response.WebResponseContent {
//	return provider.del[system.SysUser](Db(), keys)
//}

// TestService 继承 BaseService[SysUser]
type TestService struct {
	*base.BaseService[partial.TestTemplateEntity]
}

// 构造函数
// 示例：在 TestService 构造函数中设置 QueryRelativeExpression，实现自动扩展查询
func NewTestService() *TestService {
	service := &TestService{
		BaseService: base.NewBaseService[partial.TestTemplateEntity](string(initialize.DbGin)),
	}
	return service
}

// 重写分页方法
func (s *TestService) GetPageData(options request.PageDataOptions) *response.PageGridData[partial.TestTemplateEntity] {
	//// 查询前设置查询条件
	//s.BaseService.QueryRelativeExpression = func(db *gorm.DB) *gorm.DB {
	//	return db.Where("Enable = ?", 1)
	//}
	//// 统计
	//s.BaseService.SummaryExpress = func(db *gorm.DB) any {
	//	res := map[string]any{}
	//	// 方式一(推荐这种 一次查询完成)
	//	db.Select(
	//		`COUNT(Enable) as Enable,
	//				SUM(CASE WHEN Gender = 0 THEN 1 ELSE 0 END) as Gender
	//		`).
	//		Scan(&res)
	//
	//	// 方式二
	//	var enableCount int64
	//	db.Where("enable = ?", 1).
	//		Count(&enableCount)
	//	res["Enable"] = enableCount
	//
	//	// 统计 Gender=0 的总数
	//	var genderCount int64
	//	db.Where("gender = ?", 0).
	//		Count(&genderCount)
	//	res["Gender"] = genderCount
	//	return res
	//}

	// 查询后处理数据
	s.BaseService.GetPageDataOnExecuted = func(list *[]partial.TestTemplateEntity) {
		//循环处理数据
		//for i := range *list {
		//	(*list)[i].Test = fmt.Sprintf("测试数据 %d", i)
		//}
	}
	return s.BaseService.GetPageData(options)
}

// 可以选择重写方法
func (s *TestService) Add(c *gin.Context, saveModel request.SaveModel) *response.WebResponseContent {
	return s.BaseService.Add(c, saveModel)
}

// 获取当前用户信息
func (s *TestService) GetCurrentUserInfo(c *gin.Context) *response.WebResponseContent {
	data := utils.GetUserInfo(c)
	return response.Ok("", data)
}
