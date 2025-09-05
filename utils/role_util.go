package utils

func IsSuperAdmin(roleId int) bool {
	return roleId == 1
}

func IsRoleIdSuperAdmin(rolesIds []int) bool {
	contains := false
	for _, roleid := range rolesIds {
		if roleid == 1 {
			contains = true
			break
		}
	}
	return contains
}
