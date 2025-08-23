package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	systemReq "github.com/wangxin5355/vol-gin-admin-api/model/system/request"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"net/http"
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
	response.Ok(c)
}

// UpdateRolePermission
// @Tags     PermissionApi
// @Summary  更新角色菜单权限
// @Produce   application/json
// @Param    data  body      systemReq.UpdateRolePermissionReq  true  ""
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

	// 更新角色权限：先清除角色所有权限，在逐个添加菜单权限
	// 这里需要调用CasbinService
	//如果action一个没有，就不给加菜单
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户角色更新成功",
	})
}
