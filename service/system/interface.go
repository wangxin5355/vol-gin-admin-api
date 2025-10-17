package system

import (
	"github.com/wangxin5355/vol-gin-admin-api/model/dto"
)

type IPermissions interface {
	GetPermissions(roleId int) (userInter []dto.Permission, err error)
	GetPermissionsMultipleRoles(roleIds []int) (userInter []dto.Permission, err error)
}
