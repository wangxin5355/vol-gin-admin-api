package partial

import "github.com/wangxin5355/vol-gin-admin-api/model/system"

// 字典数据 (Sys_Dictionary)
type SysDictionaryEntity struct {
    system.SysDictionary
    //在这写自定义字段 例如 Test string `json:"test" gorm:"-"` //gorm:"-"表示忽略该字段
}