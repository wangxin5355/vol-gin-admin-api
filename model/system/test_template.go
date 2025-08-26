package system

import "time"

type TestTemplate struct {
	ID         uint `gorm:"primarykey"`
	Name       string
	CreateID   uint32    `gorm:"default:null"`
	CreateDate time.Time `gorm:"default:null"`
	Creator    string    `gorm:"default:null"`
	Modifier   string    `gorm:"default:null"`
	ModifyDate time.Time `gorm:"default:null"`
	ModifyID   uint32    `gorm:"default:null"`
}

func (TestTemplate) TableName() string {
	return "test_template"
}
