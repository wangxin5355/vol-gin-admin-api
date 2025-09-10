package test
        
        
import (    "time"
)

// 测试模版 (test_template)
type TestTemplate struct {
    Id int64 `json:"Id" gorm:"column:Id;primaryKey"`
    Name *string `json:"Name" gorm:"column:Name;comment:名称"`
    CreateID *int32 `json:"CreateID" gorm:"column:CreateID"`
    CreateDate *time.Time `json:"CreateDate" gorm:"column:CreateDate"`
    Creator *string `json:"Creator" gorm:"column:Creator"`
    Modifier *string `json:"Modifier" gorm:"column:Modifier"`
    ModifyDate *time.Time `json:"ModifyDate" gorm:"column:ModifyDate"`
    ModifyID *int32 `json:"ModifyID" gorm:"column:ModifyID"`
}
