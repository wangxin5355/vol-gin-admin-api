package request

// LoadTableInfoReq 获取表信息请求参数
type LoadTableInfoReq struct {
	ParentId     int    `json:"parentId"`
	TableName    string `json:"tableName"`
	ColumnCNName string `json:"columnCNName"`
	NameSpace    string `json:"nameSpace"`
	FolderName   string `json:"folderName"`
	TableId      int    `json:"tableId"`
	IsTreeLoad   bool   `json:"isTreeLoad"`
	DBServer     string `json:"dbServer"`
}
