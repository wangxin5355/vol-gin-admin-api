package system

import "time"

// SysDictionary 对应数据库表 `sys_dictionary`
type SysDictionary struct {
	DicID      int        `gorm:"column:Dic_ID;primaryKey;autoIncrement" json:"Dic_ID"`
	Config     string     `gorm:"column:Config;type:text" json:"Config"`
	CreateDate *time.Time `gorm:"column:CreateDate" json:"CreateDate"`
	CreateID   *int       `gorm:"column:CreateID" json:"CreateID"`
	Creator    *string    `gorm:"column:Creator;size:30" json:"Creator"`
	DBServer   *string    `gorm:"column:DBServer;type:text" json:"DBServer"`
	DbSql      *string    `gorm:"column:DbSql;type:text" json:"DbSql"`
	DicName    string     `gorm:"column:DicName;size:100;not null" json:"DicName"`
	DicNo      string     `gorm:"column:DicNo;size:100;not null" json:"DicNo"`
	Enable     int        `gorm:"column:Enable;not null" json:"Enable"`
	Modifier   *string    `gorm:"column:Modifier;size:30" json:"Modifier"`
	ModifyDate *time.Time `gorm:"column:ModifyDate" json:"ModifyDate"`
	ModifyID   *int       `gorm:"column:ModifyID" json:"ModifyID"`
	OrderNo    *int       `gorm:"column:OrderNo" json:"OrderNo"`
	ParentId   int        `gorm:"column:ParentId;not null" json:"ParentId"`
	Remark     *string    `gorm:"column:Remark;type:text" json:"Remark"`
}

// TableName 指定表名
func (SysDictionary) TableName() string {
	return "sys_dictionary"
}
