package request

// 更新用户角色的结构体
// UpdateUserRoleReq structure
type UpdateUserRoleReq struct {
	UserId  int    `json:"userId"`
	RoleIds string `json:"roleIds" example:"角色ID,多个逗号分隔"` // 角色ID,多个逗号分隔
}

// 更新角色菜单的结构体
// UpdateRolePermissionReq structure
type UpdateRolePermissionReq struct {
	RoleId     int          `json:"roleId"`
	MenuAction []MenuAction `json:"menuAction"` // 菜单权限数组
}

// 菜单权限
type MenuAction struct {
	MenuId  int      `json:"menuId"`                                                                // 菜单ID
	Actions []string `json:"actions" example:"Search,Add,Delete,Update,Import,Export,Upload,Audit"` // 菜单权限
}

type CheckRolePermissionReq struct {
	RoleId int    `json:"roleId"`
	MenuId int    `json:"menuId"`
	Action string `json:"action" example:"Search"` // 菜单权限
}
