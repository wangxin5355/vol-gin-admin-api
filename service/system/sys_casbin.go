package system

import (
	"fmt"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"log"
)

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

// 为用户分配多个角色
func (s *CasbinService) AssignUserRoles(userID string, roleIDs []string) error {
	// 先清除用户现有角色
	_, err := utils.GetCasbin().RemoveFilteredGroupingPolicy(0, userID)
	if err != nil {
		return err
	}
	// 添加新角色
	for _, roleID := range roleIDs {
		_, err := utils.GetCasbin().AddGroupingPolicy(userID, roleID)
		if err != nil {
			return err
		}
	}
	return utils.GetCasbin().SavePolicy()
}

// 获取用户的所有角色
func (s *CasbinService) GetUserRoles(userID string) ([]string, error) {
	return utils.GetCasbin().GetRolesForUser(userID)
}

// 检查用户权限（支持多角色）
func (s *CasbinService) CheckPermission(userID string, menuID, action string) (bool, error) {
	return utils.GetCasbin().Enforce(userID, menuID, action)
}

// 为角色添加菜单权限
func (s *CasbinService) AddMenuPermission(roleID string, menuID string, actions []string) error {
	for _, action := range actions {
		_, err := utils.GetCasbin().AddPolicy(roleID, menuID, action, "allow")
		if err != nil {
			return err
		}
	}
	return utils.GetCasbin().SavePolicy()
}

// 移除角色的菜单权限
func (s *CasbinService) RemoveMenuPermission(roleID string, menuID string, actions []string) error {
	for _, action := range actions {
		_, err := utils.GetCasbin().RemovePolicy(roleID, menuID, action, "allow")
		if err != nil {
			return err
		}
	}
	return utils.GetCasbin().SavePolicy()
}

// 获取角色对某个菜单的所有权限
func (s *CasbinService) GetRoleMenuPermissions(roleID string, menuID string) ([]string, error) {
	policies, err := utils.GetCasbin().GetFilteredPolicy(0, roleID, menuID)
	if err != nil {
		return nil, err
	}
	var actions []string
	for _, policy := range policies {
		if len(policy) > 2 && policy[3] == "allow" {
			actions = append(actions, policy[2])
		}
	}
	return actions, nil
}

// 删除角色所有菜单权限
func (s *CasbinService) RemoveMenuPermissionsByRole(roleID string) (bool, error) {
	removed, err := utils.GetCasbin().RemoveFilteredPolicy(0, roleID)
	if err != nil {
		return false, fmt.Errorf("failed to remove role menu permissions: %v", err)
	}
	log.Printf("Removed %t permissions for role %s ", removed, roleID)
	return removed, nil
}
