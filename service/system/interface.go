package system

import "github.com/wangxin5355/vol-gin-admin-api/model/system"

type IPermissions interface {
	GetPermissions(roleId int) (userInter []system.Permission, err error)
	GetPermissionsMultipleRoles(roleIds []int) (userInter []system.Permission, err error)
}
