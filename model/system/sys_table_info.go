package system

type SysTableInfo struct {
	TableId        int    `gorm:"column:Table_Id;primaryKey;autoIncrement" json:"tableId"`
	CnName         string `gorm:"column:CnName" json:"cnName"`
	ColumnCNName   string `gorm:"column:ColumnCNName" json:"columnCNName"`
	DBServer       string `gorm:"column:DBServer" json:"dbServer"`
	DataTableType  string `gorm:"column:DataTableType" json:"dataTableType"`
	DetailCnName   string `gorm:"column:DetailCnName" json:"detailCnName"`
	DetailName     string `gorm:"column:DetailName" json:"detailName"`
	EditorType     string `gorm:"column:EditorType" json:"editorType"`
	Enable         int    `gorm:"column:Enable" json:"enable"`
	ExpressField   string `gorm:"column:ExpressField" json:"expressField"`
	FolderName     string `gorm:"column:FolderName" json:"folderName"`
	Namespace      string `gorm:"column:Namespace" json:"namespace"`
	OrderNo        int    `gorm:"column:OrderNo" json:"orderNo"`
	ParentId       int    `gorm:"column:ParentId" json:"parentId"`
	RichText       string `gorm:"column:RichText" json:"richText"`
	SortName       string `gorm:"column:SortName" json:"sortName"`
	Table_Name     string `gorm:"column:TableName" json:"tableName"`
	TableTrueName  string `gorm:"column:TableTrueName" json:"tableTrueName"`
	UploadField    string `gorm:"column:UploadField" json:"uploadField"`
	UploadMaxCount int    `gorm:"column:UploadMaxCount" json:"uploadMaxCount"`
	//关联table_columns
	TableColumns []SysTableColumn `gorm:"foreignKey:TableId;references:TableId" json:"tableColumns"`
}

func (SysTableInfo) TableName() string {
	return "sys_tableinfo"
}
