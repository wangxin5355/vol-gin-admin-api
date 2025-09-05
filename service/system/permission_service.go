package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// / <summary>
// / 每个角色ID对应的菜单权限（已做静态化处理）
// / 每次获取权限时用当前服务器的版本号与redis/memory缓存的版本比较,如果不同会重新刷新缓存
// / </summary>
var rolePermissionsVersion sync.Map

// / <summary>
// / 获取角色权限时通过安全字典锁定的角色id
// / </summary>
var objKeyValue sync.Map

// / <summary>
// / 每个角色ID对应的菜单权限（已做静态化处理）
// / 每次获取权限时用当前服务器的版本号与redis/memory缓存的版本比较,如果不同会重新刷新缓存
// / </summary>
var rolePermissions sync.Map

type PermissionService struct{}

var (
	permissionServiceInstance *PermissionService
	once                      sync.Once // 确保初始化代码只执行一次
)

func GetPermissionService() *PermissionService {
	once.Do(func() {
		permissionServiceInstance = &PermissionService{}
		fmt.Println("PermissionService 单例已初始化")
	})
	return permissionServiceInstance
}

// 获取角色权限，并缓存
func (permissionService *PermissionService) GetPermissions(roleId int) (userInter []system.Permission, err error) {
	var roleIdStr = strconv.Itoa(roleId)
	var roleKey = getRoleIdKey(roleIdStr)
	cacheRolePermissionsVersion, err := global.GVA_REDIS.Get(context.Background(), roleKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		//如果redis报错，并且不是key不存在。就返回错误
		return nil, err
	}
	//角色有缓存，并且当前服务器的角色版本号与redis/memory缓存角色的版本号相同直接返回静态对象角色权限
	if currnetVeriosn, isOk := rolePermissionsVersion.Load(roleIdStr); isOk && cacheRolePermissionsVersion == currnetVeriosn.(string) {
		if permissionCahe, ok := rolePermissions.Load(roleIdStr); ok {
			if permissions, assertOk := permissionCahe.([]system.Permission); assertOk {
				return permissions, nil
			}
		}
	}
	lock, _ := objKeyValue.LoadOrStore(roleIdStr, new(sync.Mutex)) //尝试存一个roleIdStr的对象锁
	objlock, _ := lock.(sync.Mutex)
	defer objlock.Unlock()
	objlock.Lock()
	////先从缓存拿，如果有就返回
	//if permissionCahe, ok := rolePermissions.Load(roleIdStr); ok {
	//	if permissions, assertOk := permissionCahe.([]system.Permission); assertOk {
	//		return permissions, nil
	//	}
	//}
	//再从DB查询
	var permissions []system.Permission
	result := global.GVA_DB.Raw("SELECT a.Menu_Id,a.ParentId,a.TableName,a.Auth as MenuAuth, b.AuthValue as UserAuth,a.MenuType FROM `sys_menu` a left join `sys_roleauth`  b on a.Menu_Id=b.Menu_Id where b.Role_Id=? and b.AuthValue !='' ORDER BY a.ParentId", roleId).Scan(&permissions)
	if result.Error != nil {
		log.Fatal(result.Error)
		return nil, result.Error
	}
	actionToArray(permissions)
	var _version, err1 = global.GVA_REDIS.Get(context.Background(), roleKey).Result()
	if err1 != nil {
		if errors.Is(err1, redis.Nil) { //redis没有这个缓存版本，新建一个
			now := time.Now()
			_version = utils.FormatTimeMillis(now)
			//将版本号写入缓存
			global.GVA_REDIS.Set(context.Background(), roleKey, _version, 0)
		} else {
			fmt.Println("获取角色权限缓存版本号错误:", err1)
		}
	}
	//刷新当前服务器角色的权限,没有就新增，存在就覆盖
	rolePermissions.Store(roleIdStr, permissions)
	//写入当前服务器的角色最新版本号
	rolePermissionsVersion.Store(roleIdStr, _version)
	return permissions, nil

}

// 获取多角色权限，并缓存
func (permissionService *PermissionService) GetPermissionsMultipleRoles(roleIds []int) (userInter []system.Permission, err error) {
	if len(roleIds) == 0 {
		return nil, errors.New("角色数组不能为空")
	}
	if utils.IsRoleIdSuperAdmin(roleIds) { //如果是超管的，全部返回所有菜单和权限
		var permissions []system.Permission
		result := global.GVA_DB.Raw("SELECT Menu_Id,ParentId,TableName,Auth as UserAuth,MenuType FROM `sys_menu` ").Scan(&permissions)
		if result.Error != nil {
			log.Fatal(result.Error)
			return nil, result.Error
		}
		menuActionToArray(permissions)
		return permissions, nil
	}
	var roleIdStr = utils.ToListIntString(roleIds, ",")
	var roleKey = getRoleIdKey(roleIdStr)

	cacheRolePermissionsVersion, err := global.GVA_REDIS.Get(context.Background(), roleKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		//如果redis报错，并且不是key不存在。就返回错误
		return nil, err
	}
	//角色有缓存，并且当前服务器的角色版本号与redis/memory缓存角色的版本号相同直接返回静态对象角色权限
	if currnetVeriosn, isOk := rolePermissionsVersion.Load(roleIdStr); isOk && cacheRolePermissionsVersion == currnetVeriosn.(string) {
		if permissionCahe, ok := rolePermissions.Load(roleIdStr); ok {
			if permissions, assertOk := permissionCahe.([]system.Permission); assertOk {
				return permissions, nil
			}
		}
	}
	lock, _ := objKeyValue.LoadOrStore(roleIdStr, new(sync.Mutex)) //尝试存一个roleIdStr的对象锁
	objlock, _ := lock.(sync.Mutex)
	defer objlock.Unlock()
	objlock.Lock()
	////先从缓存拿，如果有就返回
	//if permissionCahe, ok := rolePermissions.Load(roleIdStr); ok {
	//	if permissions, assertOk := permissionCahe.([]system.Permission); assertOk {
	//		return permissions, nil
	//	}
	//}
	//再从DB查询
	var permissions []system.Permission
	result := global.GVA_DB.Raw("SELECT a.Menu_Id,a.ParentId,a.TableName,a.Auth as MenuAuth, b.AuthValue as UserAuth,a.MenuType FROM `sys_menu` a left join `sys_roleauth`  b on a.Menu_Id=b.Menu_Id where b.Role_Id in (?) and b.AuthValue !='' ORDER BY a.ParentId", roleIds).Scan(&permissions)
	if result.Error != nil {
		log.Fatal(result.Error)
		return nil, result.Error
	}
	actionToArray(permissions)
	var _version, err1 = global.GVA_REDIS.Get(context.Background(), roleKey).Result()
	if err1 != nil {
		if errors.Is(err1, redis.Nil) { //redis没有这个缓存版本，新建一个
			now := time.Now()
			_version = utils.FormatTimeMillis(now)
			//将版本号写入缓存
			global.GVA_REDIS.Set(context.Background(), roleKey, _version, 0)
		} else {
			fmt.Println("获取角色权限缓存版本号错误:", err1)
		}
	}
	//刷新当前服务器角色的权限,没有就新增，存在就覆盖
	rolePermissions.Store(roleIdStr, permissions)
	//写入当前服务器的角色最新版本号
	rolePermissionsVersion.Store(roleIdStr, _version)
	return permissions, nil
}

func getRoleIdKey(roleId string) string {
	return model.Role.String() + roleId
}

// 将action解析出来，并拼接
func actionToArray(permissions []system.Permission) (err error) {
	for _, permission := range permissions {
		permission.TableName = strings.ToLower(permission.TableName)
		//if permission.MenuAuth==""{
		//permission.UserAuthArr = []string{}
		//continue
		//}
		var actions []system.Action
		err := json.Unmarshal([]byte(permission.MenuAuth), &actions)
		if err != nil {
			log.Fatal("反序列化失败:", err)
			return err
		}
		//获取UserAuth和actions中都存在的权限
		permission.UserAuthArr = utils.FindCommonElementsSimple(strings.Split(permission.UserAuth, ","), extractActionValue(actions))
	}
	return nil
}

func menuActionToArray(permissions []system.Permission) (err error) {
	for _, permission := range permissions {
		permission.TableName = strings.ToLower(permission.TableName)
		//if permission.MenuAuth==""{
		//permission.UserAuthArr = []string{}
		//continue
		//}
		if permission.UserAuth == "" || permission.UserAuth == "[]" {
			continue
		}
		var actions []system.Action
		err := json.Unmarshal([]byte(permission.UserAuth), &actions)
		if err != nil {
			log.Fatal("反序列化失败:", err)
			return err
		}
		//只获取菜单的Auth里面的action
		permission.UserAuthArr = extractActionValue(actions)
	}
	return nil
}

func extractActionValue(actions []system.Action) []string {
	values := make([]string, len(actions))
	for i, action := range actions {
		values[i] = action.Value
	}
	return values
}
