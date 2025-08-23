package test

import (
	"github.com/wangxin5355/vol-gin-admin-api/core/base/provider"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
)

// TestService 通过组合“继承”ServiceBase
type TestService struct {
}

// GetPageData 获取分页数据
func (s *TestService) GetPageData(options provider.PageDataOptions) *provider.PageGridData[system.SysUser] {
	// 调用父类方法：return s.ServiceBase.GetPageData(options)
	//limit := options.Rows
	//offset := (options.Page - 1) * limit
	//db := global.GVA_DB.Model(&system.SysUser{})

	//var userLise []system.SysUser
	//options.Wheres = "[{\"name\":\"Username\",\"value\":\"u\",\"displayType\":\"like\"}]"
	//db = provider.ApplyJsonToDB(db, options)
	////db.Limit(limit).Offset(offset).Find(&userLise)
	//db.Find(&userLise)
	//return &provider.PageGridData[system.SysUser]{
	//	Rows:  userLise,
	//	Total: 100,
	//}
	return provider.GetPageData[system.SysUser](options)
}
