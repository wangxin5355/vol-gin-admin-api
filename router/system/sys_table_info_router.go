package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type SysTableInfoRouter struct{}

var sysTableInfoApi = api.ApiGroupApp.SystemApiGroup.SysTableInfoApi

func (s *SysTableInfoRouter) InitSysTableInfoRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("builder")
	{
		baseRouter.GET("getTableTree", sysTableInfoApi.GetTableTree)
		baseRouter.POST("loadTableInfo", sysTableInfoApi.LoadTableInfo)
	}
	return baseRouter
}
