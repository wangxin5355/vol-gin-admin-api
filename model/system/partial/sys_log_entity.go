package partial

import "github.com/wangxin5355/vol-gin-admin-api/model/system"

// 系统日志 (Sys_Log)
type SysLogEntity struct {
    system.SysLog
    //在这写自定义字段 例如 Test string `json:"test" gorm:"-"` //gorm:"-"表示忽略该字段
}