package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type MenuRouter struct{}

var menuApi = api.ApiGroupApp.SystemApiGroup.MenuApi

func (s *PermissionRouter) InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("menu")
	{
		baseRouter.POST("GetTreeMenu", menuApi.GetTreeMenu)
		baseRouter.GET("GetTreeMenu", menuApi.GetTreeMenu)

	}
	return baseRouter
}
