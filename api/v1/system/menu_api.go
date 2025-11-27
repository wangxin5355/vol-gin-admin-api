package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"go.uber.org/zap"
)

type MenuApi struct {
}

// GetTreeMenu
// @Tags     MenuApi
// @Summary  获取首页菜单树
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success  200   {object}  []dto.TreeMenu  "返回的菜单树json"
// @Router   /menu/GetTreeMenu [get]
func (api *MenuApi) GetTreeMenu(c *gin.Context) {
	//测试从contextt拿用户信息
	userinfo, _ := c.Get("userinfo")
	u, _ := userinfo.(*system.SysUser)
	global.GVA_LOG.Info(u.UserName)
	//---------------
	//获取当前用户roleids
	roleIds, err := utils.GetUserRolesStr(c)
	if err != nil {
		global.GVA_LOG.Error("获取用户角色失败", zap.Error(err))
		response.FailWithMessage("获取用户角色失败！无法获取菜单！", c)
		return
	}
	var menuType = utils.GetMenuType(c)
	menuTrees, err := service.ServiceInstances.MenuService.GetMenuActionList(roleIds, menuType)
	//直接返回的菜单 json 数组
	if err != nil {
		global.GVA_LOG.Error("获取菜单树失败", zap.Error(err))
		response.FailWithMessage("获取菜单树失败", c)
		return
	}
	c.JSON(http.StatusOK, menuTrees)
}

// GetMenu
// @Tags     MenuApi
// @Summary  获取所有菜单
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success  200   {object}  response.Response{data=[]system.SysMenu}  "返回所有菜单"
// @Router   /menu/GetMenu [post]
func (api *MenuApi) GetMenu(c *gin.Context) {
	res, err := service.ServiceInstances.MenuService.GetMenu()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetTreeItem
// @Tags     MenuApi
// @Summary  编辑菜单时，获取菜单信息
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param    menuId  query  int  true  "菜单ID"
// @Success  200   {object}  response.Response{data=[]system.SysMenu}  "编辑菜单时，获取菜单信息"
// @Router   /menu/GetTreeItem [get]
func (api *MenuApi) GetTreeItem(c *gin.Context) {
	menuId, err := strconv.Atoi(c.Query("menuId"))
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	res := service.ServiceInstances.MenuService.GetTreeItem(menuId)
	//为了和前端保持一致，把首字母改成小写
	resLower := make(map[string]any)
	for k, v := range res {
		resLower[utils.FirstLetterLower(k)] = v
	}
	c.JSON(http.StatusOK, resLower)
}

// Save
// @Tags     MenuApi
// @Summary  新建或编辑菜单
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param    menu  body  system.SysMenu  true  "菜单信息"
// @Success  200   {object}  response.Response{data=string}  "新建或编辑菜单"
// @Router   /menu/Save [post]
func (api *MenuApi) Save(c *gin.Context) {
	var menu system.SysMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		response.FailWithMessage("参数错误"+err.Error(), c)
		return
	}
	res := service.ServiceInstances.MenuService.Save(&menu)
	c.JSON(http.StatusOK, res)
}

// DelMenu
// @Tags     MenuApi
// @Summary  删除菜单
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param    menuId  query  int  true  "菜单ID"
// @Success  200   {object}  response.Response{data=string}  "删除菜单"
// @Router   /menu/DelMenu [post]
func (api *MenuApi) DelMenu(c *gin.Context) {
	menuId, err := strconv.Atoi(c.Query("menuId"))
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	res := service.ServiceInstances.MenuService.DelMenu(menuId)
	c.JSON(http.StatusOK, res)
}
