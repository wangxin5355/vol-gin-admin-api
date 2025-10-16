package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type AccountRouter struct{}

var accApi = api.ApiGroupApp.SystemApiGroup.AccountApi

// 注册登录路由和注册，不需要鉴权的
func (s *AccountRouter) InitAccRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("user")
	{
		baseRouter.POST("login", accApi.Login)
		baseRouter.POST("register", accApi.Register)
	}
	return baseRouter
}
