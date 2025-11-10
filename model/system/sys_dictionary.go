package system

import (
	"time"
)

// 字典数据 (Sys_Dictionary)
type SysDictionary struct {
	//第一项是固定的 写描述信息
	_          struct{}   `entity:"TableCnName=字典数据;TableName=Sys_Dictionary;DetailTable=Sys_DictionaryList;DetailTableCnName=字典明细;DBServer=gin;Key=Dic_ID"`
	ModifyDate *time.Time `json:"ModifyDate" gorm:"column:ModifyDate;comment:修改时间"`
	Modifier   *string    `json:"Modifier" gorm:"column:Modifier;comment:修改人"`
	ModifyID   int        `json:"ModifyID" gorm:"column:ModifyID"`
	CreateDate *time.Time `json:"CreateDate" gorm:"column:CreateDate;comment:创建时间"`
	Creator    *string    `json:"Creator" gorm:"column:Creator;comment:创建人"`
	CreateID   int        `json:"CreateID" gorm:"column:CreateID"`
	Enable     byte       `json:"Enable" gorm:"column:Enable;comment:是否启用"`
	Remark     *string    `json:"Remark" gorm:"column:Remark;comment:备注"`
	OrderNo    int        `json:"OrderNo" gorm:"column:OrderNo;comment:排序号"`
	DBServer   *string    `json:"DBServer" gorm:"column:DBServer;comment:DBServer"`
	DbSql      *string    `json:"DbSql" gorm:"column:DbSql;comment:sql语句"`
	Config     *string    `json:"Config" gorm:"column:Config;comment:配置项"`
	ParentId   int        `json:"ParentId" gorm:"column:ParentId;comment:父级ID"`
	DicName    string     `json:"DicName" gorm:"column:DicName;comment:字典名称"`
	DicNo      string     `json:"DicNo" gorm:"column:DicNo;comment:字典编号"`
	Dic_ID     int        `json:"Dic_ID" gorm:"column:Dic_ID;primaryKey;comment:字典ID"`
}

func (SysDictionary) TableName() string {
	return "Sys_Dictionary"
}
