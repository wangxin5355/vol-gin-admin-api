package system

import (
	"github.com/gin-gonic/gin"
	api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
)

type PermissionRouter struct{}

var permissionApi = api.ApiGroupApp.SystemApiGroup.PermissionApi

func (s *PermissionRouter) InitPermissionRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("permission")
	{
		baseRouter.POST("UpdateRolePermission", permissionApi.UpdateRolePermission)
		baseRouter.POST("UpdateUserRoles", permissionApi.UpdateUserRoles)
		baseRouter.POST("CheckRolePermission", permissionApi.CheckRolePermission)
	}
	return baseRouter
}
