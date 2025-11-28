package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type BuilderRouter struct{}

var builderApi = api.ApiGroupApp.SystemApiGroup.BuilderApi

func (s *BuilderRouter) InitBuilderRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("builder")
	{
		baseRouter.GET("GetTableTree", builderApi.GetTableTree)
		baseRouter.POST("GetTableTree", builderApi.GetTableTree)
		//baseRouter.POST("Save", builderApi.Save)
		baseRouter.POST("CreateModel", builderApi.CreateModel)
		baseRouter.POST("CreateServices", builderApi.CreateServices)
	}
	return baseRouter
}
