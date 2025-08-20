package model

import (
	"time"
)

type BaseModel struct {
	ID         uint64    `gorm:"primarykey" json:"ID"` // 主键ID
	CreateDate time.Time // 创建时间
	ModifyDate time.Time // 更新时间
}
