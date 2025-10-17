package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type SysDictionaryRouter struct{}

var sysDictionaryApi = api.ApiGroupApp.SystemApiGroup.SysDictionaryApi

func (s *SysDictionaryRouter) InitSysDictionaryRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("Sys_Dictionary")
	{
		baseRouter.GET("GetBuilderDictionary", sysDictionaryApi.GetBuilderDictionary)
	}
	return baseRouter
}
