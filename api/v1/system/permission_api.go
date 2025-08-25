package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	systemReq "github.com/wangxin5355/vol-gin-admin-api/model/system/request"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"strconv"
	"strings"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

type PermissionApi struct{}

// UpdateUserRoles
// @Tags     PermissionApi
// @Summary  更新用户角色
// @Produce   application/json
// @Param    data  body      systemReq.UpdateUserRoleReq  true  "用户名, 角色id:1,3,4"
// @Success  200   {object}  response.Response{msg=string}  "失败返回非0错误码"
// @Router   /permission/UpdateUserRoles [post]
func (api *PermissionApi) UpdateUserRoles(c *gin.Context) {
	var req systemReq.UpdateUserRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.UpdateUserRoleVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 更新用户角色
	// 这里需要调用CasbinService的AssignUserRoles方法
	strArray := strings.Split(req.RoleIds, ",")
	err = casbinService.AssignUserRoles(strconv.Itoa(req.UserId), strArray)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok("更新用户角色成功", c)
}

// UpdateRolePermission
// @Tags     PermissionApi
// @Summary  更新角色菜单权限
// @Produce   application/json
// @Param    data  body      systemReq.UpdateRolePermissionReq  true  " "
// @Success  200   {object}  response.Response{msg=string}  "失败返回非0错误码"
// @Router   /permission/UpdateRolePermission [post]
func (api *PermissionApi) UpdateRolePermission(c *gin.Context) {
	var req systemReq.UpdateRolePermissionReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//err = utils.Verify(req, utils.UpdateRolePermissionVerify)
	//if err != nil {
	//response.FailWithMessage(err.Error(), c)
	//return
	//}
	//逻辑代码写校验
	if req.RoleId <= 0 {
		response.FailWithMessage("参数RoleId异常", c)
		return
	}
	if req.MenuAction == nil || len(req.MenuAction) == 0 {
		response.FailWithMessage("菜单信息不能为空", c)
		return
	}

	// 更新角色权限：先清除角色所有权限，在逐个添加菜单权限
	// 这里需要调用CasbinService
	_, err = casbinService.RemoveMenuPermissionsByRole(strconv.Itoa(req.RoleId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, menuAction := range req.MenuAction {
		if menuAction.MenuId <= 0 {
			continue
		}
		if menuAction.Actions == nil || len(menuAction.Actions) == 0 {
			continue
		}
		err = casbinService.AddMenuPermission(strconv.Itoa(req.RoleId), strconv.Itoa(menuAction.MenuId), menuAction.Actions)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	response.Ok("更新角色权限成功", c)
}

// CheckRolePermission
// @Tags     PermissionApi
// @Summary  检查角色权限
// @Produce   application/json
// @Param    data  body      systemReq.CheckRolePermissionReq  true  " "
// @Success  200   {object}  response.Response{data=bool,msg=string}  "失败返回非0错误码"
// @Router   /permission/CheckRolePermission [post]
func (api *PermissionApi) CheckRolePermission(c *gin.Context) {
	var req systemReq.CheckRolePermissionReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.CheckRolePermissionVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	hasPermission, err1 := casbinService.CheckPermission(strconv.Itoa(req.RoleId), strconv.Itoa(req.MenuId), req.Action)
	if err1 != nil {
		response.FailWithMessage(err1.Error(), c)
		return
	}
	response.OkWithData(hasPermission, c)

}
