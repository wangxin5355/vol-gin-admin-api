package response

import (
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
)

type SysUserResp struct {
	User system.SysUser `json:"user"`
}

type LoginResp struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}
