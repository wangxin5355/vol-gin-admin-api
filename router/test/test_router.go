package test

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type TestRouter struct{}

var testApi = api.ApiGroupApp.TestApiGroup.TestApi

func (s *TestRouter) InitTestRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("test")
	{
		baseRouter.POST("GetPageData", testApi.GetPageData)
	}
	return baseRouter
}
