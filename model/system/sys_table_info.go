package system

type SysTableInfo struct {
	TableId        int    `gorm:"column:Table_Id;primaryKey;autoIncrement" json:"tableId"`
	CnName         string `gorm:"column:CnName;default:null" json:"cnName"`
	ColumnCNName   string `gorm:"column:ColumnCNName;default:null" json:"columnCNName"`
	DBServer       string `gorm:"column:DBServer;default:null" json:"dbServer"`
	DataTableType  string `gorm:"column:DataTableType;default:null" json:"dataTableType"`
	DetailCnName   string `gorm:"column:DetailCnName;default:null" json:"detailCnName"`
	DetailName     string `gorm:"column:DetailName;default:null" json:"detailName"`
	EditorType     string `gorm:"column:EditorType;default:null" json:"editorType"`
	Enable         int    `gorm:"column:Enable;default:null" json:"enable"`
	ExpressField   string `gorm:"column:ExpressField;default:null" json:"expressField"`
	FolderName     string `gorm:"column:FolderName;default:null" json:"folderName"`
	Namespace      string `gorm:"column:Namespace;default:null" json:"namespace"`
	OrderNo        int    `gorm:"column:OrderNo;default:null" json:"orderNo"`
	ParentId       int    `gorm:"column:ParentId;default:null" json:"parentId"`
	RichText       string `gorm:"column:RichText;default:null" json:"richText"`
	SortName       string `gorm:"column:SortName;default:null" json:"sortName"`
	Table_Name     string `gorm:"column:TableName;default:null" json:"tableName"`
	TableTrueName  string `gorm:"column:TableTrueName;default:null" json:"tableTrueName"`
	UploadField    string `gorm:"column:UploadField;default:null" json:"uploadField"`
	UploadMaxCount int    `gorm:"column:UploadMaxCount;default:null" json:"uploadMaxCount"`
	//关联table_columns
	TableColumns []SysTableColumn `gorm:"foreignKey:TableId;references:TableId" json:"tableColumns"`
}

func (SysTableInfo) TableName() string {
	return "Sys_TableInfo"
}
