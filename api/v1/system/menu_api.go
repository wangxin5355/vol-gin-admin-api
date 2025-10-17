package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"go.uber.org/zap"
	"net/http"
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
