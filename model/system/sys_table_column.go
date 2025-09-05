package system

import "time"

type SysTableColumn struct {
	ColumnId      int       `gorm:"primaryKey;column:ColumnId" json:"columnId"`
	ApiInPut      int       `gorm:"column:ApiInPut;default:null" json:"apiInPut"`
	ApiIsNull     int       `gorm:"column:ApiIsNull;default:null" json:"apiIsNull"`
	ApiOutPut     int       `gorm:"column:ApiOutPut;default:null" json:"apiOutPut"`
	ColSize       int       `gorm:"column:ColSize;default:null" json:"colSize"`
	ColumnCNName  string    `gorm:"column:ColumnCNName;default:null" json:"columnCNName"`
	ColumnName    string    `gorm:"column:ColumnName;default:null" json:"columnName"`
	ColumnType    string    `gorm:"column:ColumnType;default:null" json:"columnType"`
	ColumnWidth   int       `gorm:"column:ColumnWidth;default:null" json:"columnWidth"`
	Columnformat  string    `gorm:"column:Columnformat;default:null" json:"columnformat"`
	CreateDate    time.Time `gorm:"column:CreateDate;default:null" json:"createDate"`
	CreateID      int       `gorm:"column:CreateID;default:null" json:"createID"`
	Creator       string    `gorm:"column:Creator;default:null" json:"creator"`
	DropNo        string    `gorm:"column:DropNo;default:null" json:"dropNo"`
	EditColNo     int       `gorm:"column:EditColNo;default:null" json:"editColNo"`
	EditRowNo     int       `gorm:"column:EditRowNo;default:null" json:"editRowNo"`
	EditType      string    `gorm:"column:EditType;default:null" json:"editType"`
	Enable        int       `gorm:"column:Enable;default:null" json:"enable"`
	IsColumnData  int       `gorm:"column:IsColumnData;default:null" json:"isColumnData"`
	IsDisplay     int       `gorm:"column:IsDisplay;default:null" json:"isDisplay"`
	IsImage       int       `gorm:"column:IsImage;default:null" json:"isImage"`
	IsKey         int       `gorm:"column:IsKey;default:null" json:"isKey"`
	IsNull        int       `gorm:"column:IsNull;default:null" json:"isNull"`
	IsReadDataset int       `gorm:"column:IsReadDataset;default:null" json:"isReadDataset"`
	Maxlength     int       `gorm:"column:Maxlength;default:null" json:"maxlength"`
	Modifier      string    `gorm:"column:Modifier;default:null" json:"modifier"`
	ModifyDate    time.Time `gorm:"column:ModifyDate;default:null" json:"modifyDate"`
	ModifyID      int       `gorm:"column:ModifyID;default:null" json:"modifyID"`
	OrderNo       int       `gorm:"column:OrderNo;default:null" json:"orderNo"`
	Script        string    `gorm:"column:Script;default:null" json:"script"`
	SearchColNo   int       `gorm:"column:SearchColNo;default:null" json:"searchColNo"`
	SearchRowNo   int       `gorm:"column:SearchRowNo;default:null" json:"searchRowNo"`
	SearchType    string    `gorm:"column:SearchType;default:null" json:"searchType"`
	Sortable      int       `gorm:"column:Sortable;default:null" json:"sortable"`
	Table_Name    string    `gorm:"column:TableName;default:null" json:"tableName"`
	TableId       int       `gorm:"column:Table_Id;default:null" json:"tableId"`
}

func (SysTableColumn) TableName() string {
	return "sys_tablecolumn"
}
