package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"sync"
	"time"
)

type MenuService struct {
}

var permissionService = GetPermissionService()

const _menuCacheKey = "inernalMenu"

var _menuVersionn = ""
var _menus []system.SysMenu
var allmenuLock sync.Mutex

// 获取对应角色的菜单列表
func (menuService *MenuService) GetMenuActionList(roleIds []string, menuType int) (treeMenus []system.TreeMenu, err error) {
	var roleIds_int = utils.StringSliceToIntSliceFilter(roleIds)
	if utils.IsRoleIdSuperAdmin(roleIds_int) { //如果是超管的，全部返回所有菜单和权限
		sys_menus, err1 := menuService.getAllMenu()
		if err1 != nil {
			return nil, err1
		}
		//转换成TreeMenu,预分配容量，提高性能
		treeMenus = make([]system.TreeMenu, 0, len(sys_menus))
		for _, sysMenu := range sys_menus {
			ps := make([]string, len(sysMenu.Actions))
			for _, p := range sysMenu.Actions {
				ps = append(ps, p.Value)
			}
			treeMenus = append(treeMenus, system.TreeMenu{
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
	permissionsIndex := make(map[int]*system.Permission)
	for i := range permissions {
		permissionsIndex[permissions[i].Menu_Id] = &permissions[i]
	}
	treeMenusIndex := make(map[int]system.TreeMenu)
	for _, sysMenu := range sys_menus {
		if sysMenu.MenuType != menuType {
			continue
		} //匹配菜单类型，pc和小程序不一样
		//从权限中找是否有次菜单，并且检查是否已经存在于用户授权，因为一个人可能两个角色，两个角色可能授权相同菜单
		if permission, exists := permissionsIndex[sysMenu.Menu_Id]; exists {
			//检查这个菜单是否已经存在于菜单树中，去重
			if _, exists := treeMenusIndex[sysMenu.Menu_Id]; !exists {
				//检查这个菜单是否已经存在于菜单树中，去重,不存在才加入
				treeMenusIndex[sysMenu.Menu_Id] = system.TreeMenu{
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
	treeMenus = make([]system.TreeMenu, 0, len(treeMenusIndex))
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
			menu.Actions = []system.Action{} //给个空切片，避免序列化问题
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
