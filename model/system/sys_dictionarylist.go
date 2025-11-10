package system

import "time"

type SysDictionaryList struct {
	DicList_ID int        `gorm:"column:DicList_ID;primaryKey;autoIncrement" json:"DicList_ID"`
	CreateDate *time.Time `gorm:"column:CreateDate" json:"CreateDate"`
	CreateID   *int       `gorm:"column:CreateID" json:"CreateID"`
	Creator    *string    `gorm:"column:Creator" json:"Creator"`
	DicName    *string    `gorm:"column:DicName" json:"DicName"`
	DicValue   *string    `gorm:"column:DicValue" json:"DicValue"`
	Dic_ID     int        `gorm:"column:Dic_ID" json:"Dic_ID"`
	Enable     *int       `gorm:"column:Enable" json:"Enable"`
	Modifier   *string    `gorm:"column:Modifier" json:"Modifier"`
	ModifyDate *time.Time `gorm:"column:ModifyDate" json:"ModifyDate"`
	ModifyID   *int       `gorm:"column:ModifyID" json:"ModifyID"`
	OrderNo    *int       `gorm:"column:OrderNo" json:"OrderNo"`
	Remark     *string    `gorm:"column:Remark" json:"Remark"`
}

func (SysDictionaryList) TableName() string {
	return "sys_dictionarylist"
}
