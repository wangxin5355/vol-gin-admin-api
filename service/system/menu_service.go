package system

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/dto"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"gorm.io/gorm"
)

type MenuService struct {
}

var permissionService = GetPermissionService()

const _menuCacheKey = "inernalMenu"

var _menuVersionn = ""
var _menus []system.SysMenu
var allmenuLock sync.Mutex

func SysMenuDB() *gorm.DB {
	return global.GVA_DB.Model(&system.SysMenu{})
}

// 获取对应角色的菜单列表
func (menuService *MenuService) GetMenuActionList(roleIds []string, menuType int) (treeMenus []dto.TreeMenu, err error) {
	var roleIds_int = utils.StringSliceToIntSliceFilter(roleIds)
	if utils.IsRoleIdSuperAdmin(roleIds_int) { //如果是超管的，全部返回所有菜单和权限
		sys_menus, err1 := menuService.getAllMenu()
		if err1 != nil {
			return nil, err1
		}
		//转换成TreeMenu,预分配容量，提高性能
		treeMenus = make([]dto.TreeMenu, 0, len(sys_menus))

		type AuthItem struct {
			Text  string `json:"text"`
			Value string `json:"value"`
		}

		for _, sysMenu := range sys_menus {
			var auths []AuthItem
			if err := json.Unmarshal([]byte(sysMenu.Auth), &auths); err != nil {
				continue
			}

			var ps []string
			for _, a := range auths {
				ps = append(ps, a.Value)
			}
			treeMenus = append(treeMenus, dto.TreeMenu{
				ID:         sysMenu.Menu_Id,
				Name:       sysMenu.MenuName,
				Url:        sysMenu.Url,
				ParentId:   sysMenu.ParentId,
				Icon:       sysMenu.Icon,
				Enable:     sysMenu.Enable,
				TableName:  sysMenu.ITableName,
				Permission: ps,
			})
		}
		return treeMenus, nil
	}
	//根据授权，根据授权表+菜单表结合获取到
	permissions, err2 := permissionService.GetPermissionsMultipleRoles(roleIds_int)
	if err2 != nil {
		return nil, err2
	}
	sys_menus, err3 := menuService.getAllMenu()
	if err3 != nil {
		return nil, err3
	}
	//从所有菜单中，找到用户授权菜单，去除，
	//预先建map，提高性能，
	permissionsIndex := make(map[int]*dto.Permission)
	for i := range permissions {
		permissionsIndex[permissions[i].Menu_Id] = &permissions[i]
	}
	treeMenusIndex := make(map[int]dto.TreeMenu)
	for _, sysMenu := range sys_menus {
		if sysMenu.MenuType != menuType {
			continue
		} //匹配菜单类型，pc和小程序不一样
		//从权限中找是否有次菜单，并且检查是否已经存在于用户授权，因为一个人可能两个角色，两个角色可能授权相同菜单
		if permission, exists := permissionsIndex[sysMenu.Menu_Id]; exists {
			//检查这个菜单是否已经存在于菜单树中，去重
			if _, exists := treeMenusIndex[sysMenu.Menu_Id]; !exists {
				//检查这个菜单是否已经存在于菜单树中，去重,不存在才加入
				treeMenusIndex[sysMenu.Menu_Id] = dto.TreeMenu{
					ID:         sysMenu.Menu_Id,
					Name:       sysMenu.MenuName,
					Url:        sysMenu.Url,
					ParentId:   sysMenu.ParentId,
					Icon:       sysMenu.Icon,
					Enable:     sysMenu.Enable,
					TableName:  sysMenu.ITableName,
					Permission: permission.UserAuthArr,
				}
			}
		}
	}
	//map转切片
	treeMenus = make([]dto.TreeMenu, 0, len(treeMenusIndex))
	for _, treeMenu := range treeMenusIndex {
		treeMenus = append(treeMenus, treeMenu)
	}
	return treeMenus, nil
}

// 获取所有菜单，并缓存到本地
func (menuService *MenuService) getAllMenu() (menus []system.SysMenu, err error) {
	//每次比较缓存是否更新过，如果更新则重新获取数据
	_cacheVersion, err := global.GVA_REDIS.Get(context.Background(), _menuCacheKey).Result()
	if _menuVersionn != "" && err == nil && _cacheVersion == _menuVersionn {
		//返回本地缓存
		return _menus, nil
	}
	defer allmenuLock.Unlock()
	allmenuLock.Lock()
	if _menuVersionn != "" && len(_menus) != 0 && _menuVersionn == _cacheVersion {
		return _menus, nil
	}
	//从DB获取所有的菜单
	err = global.GVA_DB.Where("Enable = 1 or Enable=2").Order("OrderNo").Order("ParentId desc").Find(&_menus).Error
	if err != nil {
		fmt.Println("获取所有菜单失败:", err)
		return _menus, err
	}
	for _, menu := range _menus {
		if menu.Auth != "" && len(menu.Auth) > 10 {
			json.Unmarshal([]byte(menu.Auth), &menu.Actions)
		} else {
			menu.Actions = []dto.Action{} //给个空切片，避免序列化问题
		}
	}
	_cacheVersion, err = global.GVA_REDIS.Get(context.Background(), _menuCacheKey).Result()
	if _cacheVersion == "" { //不管是没拿到，还是redis出错，这里如果只要是"".就重建缓存
		now := time.Now()
		_cacheVersion = utils.FormatTimeMillis(now)
		err = global.GVA_REDIS.Set(context.Background(), _menuCacheKey, _cacheVersion, 0).Err()
		if err != nil {
			return nil, err
		}
		_menuVersionn = _cacheVersion
	} else {
		_menuVersionn = _cacheVersion
	}
	return _menus, nil
}

// GetMenu 获取所有菜单列表
func (menuService *MenuService) GetMenu() ([]map[string]any, error) {
	var menus []system.SysMenu
	SysMenuDB().Order("OrderNo").Order("ParentId").Find(&menus)
	var res []map[string]any
	for _, m := range menus {
		res = append(res, map[string]any{
			"id":       m.Menu_Id,
			"parentId": m.ParentId,
			"name":     m.MenuName,
			"icon":     m.Icon,
			"menuType": m.MenuType,
			"orderNo":  m.OrderNo,
		})
	}
	return res, nil
}

// GetTreeItem 编辑菜单时，获取菜单信息
func (menuService *MenuService) GetTreeItem(menuId int) map[string]any {
	var res map[string]any
	err := SysMenuDB().
		Where("Menu_Id = ?", menuId).
		Select(`Menu_Id, ParentId, MenuName, Url, Auth, OrderNo, Icon, Enable,
            COALESCE(MenuType, 0) as MenuType, CreateDate, Creator, TableName, ModifyDate`).
		Scan(&res).Error
	if err != nil {
		return nil
	}
	return res
}

// Save 新建或编辑菜单
func (menuService *MenuService) Save(c *gin.Context) *response.WebResponseContent {
	var menu system.SysMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		return response.Error("参数错误" + err.Error())
	}
	if menu.Menu_Id > 0 && menu.Menu_Id == menu.ParentId {
		return response.Error("父ID不能和菜单ID相同")
	}
	if utils.IsNull(menu.MenuName) || utils.IsNull(menu.ITableName) {
		return response.Error("菜单名称和表名称不能为空")
	}
	if menu.ITableName != "/" && menu.ITableName != "." {
		//判断一下是否存在
		sysMenu := &system.SysMenu{}
		err := global.GVA_DB.Where("TableName = ?", menu.ITableName).First(sysMenu).Error
		if err != nil {
			return response.Error(err.Error())
		}
		if sysMenu != nil {
			if sysMenu.MenuType == menu.MenuType {
				if (menu.Menu_Id > 0 && sysMenu.Menu_Id != menu.Menu_Id) || (menu.Menu_Id <= 0) {
					return response.Error("表/视图【" + menu.ITableName + "】名称已被其他菜单使用")
				}
			}
		}
		changed := false
		//TODO: 添加和修改需要做默认值 这个默认值需要重新封装一下
		if menu.Menu_Id <= 0 {
			//新增
			err := SysMenuDB().Create(menu).Error
			if err != nil {
				return response.Error("新增菜单失败：" + err.Error())
			}
		} else {
			//编辑
			if menu.Menu_Id == menu.ParentId {
				return response.Error("父ID不能和菜单ID相同")
			}
			if (SysMenuDB().Where("ParentId = ? AND Menu_Id = ?", menu.Menu_Id, menu.ParentId).First(&system.SysMenu{}).RowsAffected > 0) {
				return response.Error("不能选择此父级id，选择的父级id与当前菜单形成依赖关系")
			}
			var auth string
			changed = SysMenuDB().Where("Menu_Id = ?", menu.Menu_Id).Select("Auth").Scan(&auth).Error == nil && auth != menu.Auth
			err := SysMenuDB().Where("Menu_Id = ?", menu.Menu_Id).
				Select("ParentId", "MenuName", "Url", "Auth", "OrderNo",
					"Icon", "Enable", "MenuType", "TableName",
					"ModifyDate", "Modifier").
				Updates(menu).Error
			if err != nil {
				return response.Error("修改菜单失败：" + err.Error())
			}
		}
		global.GVA_REDIS.Set(context.Background(), _menuCacheKey, utils.FormatTimeMillis(time.Now()), 0)
		//TODO: 如果权限有变化，则需要更新角色的权限缓存
		if changed == true {
			//要更新角色的权限缓存
		}
	}
	return response.Ok("操作成功", menu)
}

// DelMenu 删除菜单
func (menuService *MenuService) DelMenu(menuId int) *response.WebResponseContent {
	if menuId <= 0 {
		return response.Error("参数错误")
	}
	//检查是否有子菜单
	if SysMenuDB().Where("ParentId = ?", menuId).First(&system.SysMenu{}).RowsAffected > 0 {
		return response.Error("请先删除子菜单")
	}
	//删除菜单
	err := SysMenuDB().Where("Menu_Id = ?", menuId).Delete(&system.SysMenu{}).Error
	if err != nil {
		return response.Error("删除失败：" + err.Error())
	}
	//更新缓存
	global.GVA_REDIS.Set(context.Background(), _menuCacheKey, utils.FormatTimeMillis(time.Now()), 0)
	return response.Ok("删除成功", nil)
}
