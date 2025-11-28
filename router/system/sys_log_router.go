package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type SysLogRouter struct{}

var SysLogApi = api.ApiGroupApp.SystemApiGroup.SysLogApi

func (s *SysLogRouter) InitSysLogRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("Sys_Log")
	{
		baseRouter.POST("getPageData", SysLogApi.GetPageData)
		baseRouter.POST("Add", SysLogApi.Add)
		baseRouter.POST("Update", SysLogApi.Update)
		baseRouter.POST("Del", SysLogApi.Del)
	}
	return baseRouter
}
