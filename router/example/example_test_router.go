package example

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
	"github.com/wangxin5355/vol-gin-admin-api/middleware"
)

type ExampleTestRouter struct{}

var exampleApi = api.ApiGroupApp.ExampleApiGroup.ExampleTestApi

func (e *ExampleTestRouter) InitTestRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("example").Use(middleware.OperationRecord()) //带记录请求日志的
	{
		customerRouter.GET("Test", exampleApi.Test)

	}
	customerRouterWithoutRecord := Router.Group("example") //不带请求日志
	{
		customerRouterWithoutRecord.GET("Test1", exampleApi.Test1)

	}
}
