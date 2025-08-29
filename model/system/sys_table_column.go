package system

import "time"

type SysTableColumn struct {
	ColumnId      int        `gorm:"primaryKey;column:ColumnId" json:"columnId"`
	ApiInPut      int        `gorm:"column:ApiInPut" json:"apiInPut"`
	ApiIsNull     int        `gorm:"column:ApiIsNull" json:"apiIsNull"`
	ApiOutPut     int        `gorm:"column:ApiOutPut" json:"apiOutPut"`
	ColSize       int        `gorm:"column:ColSize" json:"colSize"`
	ColumnCNName  string     `gorm:"column:ColumnCNName" json:"columnCNName"`
	ColumnName    string     `gorm:"column:ColumnName" json:"columnName"`
	ColumnType    string     `gorm:"column:ColumnType" json:"columnType"`
	ColumnWidth   int        `gorm:"column:ColumnWidth" json:"columnWidth"`
	Columnformat  string     `gorm:"column:Columnformat" json:"columnformat"`
	CreateDate    *time.Time `gorm:"column:CreateDate" json:"createDate"`
	CreateID      int        `gorm:"column:CreateID" json:"createID"`
	Creator       string     `gorm:"column:Creator" json:"creator"`
	DropNo        string     `gorm:"column:DropNo" json:"dropNo"`
	EditColNo     int        `gorm:"column:EditColNo" json:"editColNo"`
	EditRowNo     int        `gorm:"column:EditRowNo" json:"editRowNo"`
	EditType      string     `gorm:"column:EditType" json:"editType"`
	Enable        int        `gorm:"column:Enable" json:"enable"`
	IsColumnData  int        `gorm:"column:IsColumnData" json:"isColumnData"`
	IsDisplay     int        `gorm:"column:IsDisplay" json:"isDisplay"`
	IsImage       int        `gorm:"column:IsImage" json:"isImage"`
	IsKey         int        `gorm:"column:IsKey" json:"isKey"`
	IsNull        int        `gorm:"column:IsNull" json:"isNull"`
	IsReadDataset int        `gorm:"column:IsReadDataset" json:"isReadDataset"`
	Maxlength     int        `gorm:"column:Maxlength" json:"maxlength"`
	Modifier      string     `gorm:"column:Modifier" json:"modifier"`
	ModifyDate    *time.Time `gorm:"column:ModifyDate" json:"modifyDate"`
	ModifyID      int        `gorm:"column:ModifyID" json:"modifyID"`
	OrderNo       int        `gorm:"column:OrderNo" json:"orderNo"`
	Script        string     `gorm:"column:Script" json:"script"`
	SearchColNo   int        `gorm:"column:SearchColNo" json:"searchColNo"`
	SearchRowNo   int        `gorm:"column:SearchRowNo" json:"searchRowNo"`
	SearchType    string     `gorm:"column:SearchType" json:"searchType"`
	Sortable      int        `gorm:"column:Sortable" json:"sortable"`
	Table_Name    string     `gorm:"column:TableName" json:"tableName"`
	TableId       int        `gorm:"column:Table_Id" json:"tableId"`
}

func (SysTableColumn) TableName() string {
	return "sys_tablecolumn"
}
