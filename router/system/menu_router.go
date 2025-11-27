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
		baseRouter.GET("getTreeMenu", menuApi.GetTreeMenu)
		baseRouter.POST("getTreeMenu", menuApi.GetTreeMenu)
		baseRouter.POST("getMenu", menuApi.GetMenu)
		baseRouter.POST("getTreeItem", menuApi.GetTreeItem)
		baseRouter.POST("save", menuApi.Save)
		baseRouter.POST("delMenu", menuApi.DelMenu)
	}
	return baseRouter
}
