package partial

import "github.com/wangxin5355/vol-gin-admin-api/model/system"

type TestTemplateEntity struct {
	system.TestTemplate
	Test string `json:"test" gorm:"-"`
	//ID         uint `gorm:"primarykey"`
	//Name       string
	//CreateID   uint32    `gorm:"default:null"`
	//CreateDate time.Time `gorm:"default:null"`
	//Creator    string    `gorm:"default:null"`
	//Modifier   string    `gorm:"default:null"`
	//ModifyDate time.Time `gorm:"default:null"`
	//ModifyID   uint32    `gorm:"default:null"`
}
