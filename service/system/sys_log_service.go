package system

import (
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"

	"github.com/wangxin5355/vol-gin-admin-api/core/base"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/model/system/partial"
)

func InitSysLogService() *SysLogService {
	return &SysLogService{
		BaseService: base.InitBaseService[partial.SysLogEntity, system.SysLog](string(initialize.DbGin)),
	}
}

// SysLogService 继承 BaseService[SysLog]
type SysLogService struct {
	*base.BaseService[partial.SysLogEntity, system.SysLog]
}
