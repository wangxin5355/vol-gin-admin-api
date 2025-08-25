package test

import (
	"github.com/wangxin5355/vol-gin-admin-api/core/base"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
)

//
//type TestService struct {
//}
//
//func Db() *gorm.DB {
//	return global.GetGlobalDBByDBName(string(initialize.DbGin))
//}
//
//// GetPageData 获取分页数据
//func (s *TestService) GetPageData(options request.PageDataOptions) *response.PageGridData[system.SysUser] {
//	// 调用父类方法：return s.ServiceBase.GetPageData(options)
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
//	return provider.GetPageData[system.SysUser](Db(), options)
//}
//
//// Add 添加方法
//func (s *TestService) Add(saveModel request.SaveModel) *response.WebResponseContent {
//	return provider.Add[system.SysUser](Db(), saveModel)
//}
//
//// Update 更新方法
//func (s *TestService) Update(saveModel request.SaveModel) *response.WebResponseContent {
//	return provider.Update[system.SysUser](Db(), saveModel)
//}
//
//// Del 删除方法
//func (s *TestService) Del(keys []any) *response.WebResponseContent {
//	return provider.Del[system.SysUser](Db(), keys)
//}

// TestService 继承 BaseService[SysUser]
type TestService struct {
	*base.BaseService[system.SysUser]
}

// 构造函数
func NewTestService() *TestService {
	return &TestService{
		BaseService: base.NewBaseService[system.SysUser](string(initialize.DbGin)),
	}
}

// 重写分页方法
func (s *TestService) GetPageData(options request.PageDataOptions) *response.PageGridData[system.SysUser] {
	return s.BaseService.GetPageData(options)
}

// 可以选择重写方法
func (s *TestService) Add(saveModel request.SaveModel) *response.WebResponseContent {
	return s.BaseService.Add(saveModel)
}
