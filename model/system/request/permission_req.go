package request

import (
	"github.com/wangxin5355/vol-gin-admin-api/model/dto"
)

// 更新用户角色的结构体
// UpdateUserRoleReq structure
type UpdateUserRoleReq struct {
	UserId  int    `json:"userId"`
	RoleIds string `json:"roleIds" example:"角色ID,多个逗号分隔"` // 角色ID,多个逗号分隔
}

// 更新角色菜单的结构体
// UpdateRolePermissionReq structure
type UpdateRolePermissionReq struct {
	RoleId     int              `json:"roleId"`
	MenuAction []dto.MenuAction `json:"menuAction"` // 菜单权限数组
}

type CheckRolePermissionReq struct {
	RoleId int    `json:"roleId"`
	MenuId int    `json:"menuId"`
	Action string `json:"action" example:"Search"` // 菜单权限
}
