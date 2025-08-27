package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/core/base"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/model/system/partial"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

// TestService 继承 BaseService[SysUser]
type TestService struct {
	*base.BaseService[partial.TestTemplateEntity, system.TestTemplate]
}

// 构造函数
// 示例：在 TestService 构造函数中设置 QueryRelativeExpression，实现自动扩展查询
func NewTestService() *TestService {
	service := &TestService{
		//这里写了两个实体，为了兼容一些扩展字段，如果只写一个转换会很麻烦
		BaseService: base.NewBaseService[partial.TestTemplateEntity, system.TestTemplate](string(initialize.DbGin)),
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
	// 保存前操作
	s.AddOnExecuting = func(entity *system.TestTemplate) *response.WebResponseContent {
		// 获取当前用户信息
		entity.Name = "保存前操作"
		//return response.Ok("", entity)
		return response.Error("保存前操作错误")
	}
	s.AddOnExecuted = func(entity *system.TestTemplate) *response.WebResponseContent {
		//return response.Ok("保存后操作", entity)
		return response.Error("保存后操作错误")
	}
	return s.BaseService.Add(c, saveModel)
}

// 重写Update
func (s *TestService) Update(c *gin.Context, saveModel request.SaveModel) *response.WebResponseContent {
	// 更新前操作
	s.UpdateOnExecuting = func(entity *system.TestTemplate) *response.WebResponseContent {
		entity.Name = "更新前操作"
		return response.Ok("", entity)
		//return response.Error("更新前操作错误")
	}
	//s.UpdateOnExecuted = func(entity *system.TestTemplate) *response.WebResponseContent {
	//	return response.Ok("更新后操作", entity)
	//	//return response.Error("更新后操作错误")
	//}
	return s.BaseService.Update(c, saveModel)
}

// 重写Del
func (s *TestService) Del(c *gin.Context, keys []any) *response.WebResponseContent {
	//// 删除前操作
	//s.DelOnExecuting = func(keys []any) *response.WebResponseContent {
	//	//循环输出keys
	//	for _, key := range keys {
	//		fmt.Println("删除前操作 key:", key)
	//	}
	//	return response.Error("删除前操作")
	//}
	//// 删除后操作
	//s.DelOnExecuted = func(keys []any) *response.WebResponseContent {
	//	//循环输出keys
	//	for _, key := range keys {
	//		println("删除后操作 key:", key)
	//	}
	//	return response.Error("删除后操作")
	//}
	return s.BaseService.Del(c, keys)
}

// 获取当前用户信息
func (s *TestService) GetCurrentUserInfo(c *gin.Context) *response.WebResponseContent {
	data := utils.GetUserInfo(c)
	return response.Ok("", data)
}
