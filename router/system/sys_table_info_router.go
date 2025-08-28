package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/api/v1/system"
)

// Api
type apiGroup struct {
	system.SysTableInfoApi
}

var apiGroupApp = new(apiGroup)

// Router
// 只保留一个 routerGroup 结构体，并包含 apiGroup 字段
// 移除重复和无用结构体

type routerGroup struct {
	apiGroup *apiGroup
}

var groupApp = &routerGroup{apiGroup: apiGroupApp}

// 修复 InitSysTableInfoRouter 方法，使用 groupApp.apiGroup
func (s *routerGroup) InitSysTableInfoRouter(Router gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("sysTableInfo")
	{
		baseRouter.GET("getTableTree", s.apiGroup.GetTableTree)
	}
	return baseRouter
}
