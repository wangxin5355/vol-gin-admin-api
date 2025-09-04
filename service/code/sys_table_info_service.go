package code

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/model/system/request"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var SysTableInfoGroup struct {
	SysTableInfoService *SysTableInfoService
}

type SysTableInfoService struct {
}

// GetTableTree 获取表结构树形数据
func (s *SysTableInfoService) GetTableTree() map[string]interface{} {
	type TreeNode struct {
		Id       int    `json:"id"`
		PId      int    `json:"pId"`
		ParentId int    `json:"parentId"`
		Name     string `json:"name"`
		OrderNo  int    `json:"orderNo"`
		IsParent bool   `json:"isParent"`
	}

	var tableInfos []system.SysTableInfo
	if err := global.GVA_DB.Order("OrderNo").Find(&tableInfos).Error; err != nil {
		return map[string]interface{}{}
	}

	// 构造所有节点的 pId 集合
	pIdSet := make(map[int]struct{})
	for _, info := range tableInfos {
		pIdSet[info.ParentId] = struct{}{}
	}

	treeList := make([]TreeNode, 0, len(tableInfos))
	for _, info := range tableInfos {
		node := TreeNode{
			Id:       info.TableId,
			PId:      info.ParentId,
			ParentId: info.ParentId,
			Name:     info.ColumnCNName,
			OrderNo:  info.OrderNo,
			IsParent: false,
		}
		if _, ok := pIdSet[info.TableId]; ok {
			node.IsParent = true
		}
		treeList = append(treeList, node)
	}
	//TODO:没有类库之分 先返回固定的，看看后边要不要按照文件夹划分
	var nameSpace []string
	nameSpace = append(nameSpace, "--")
	res := map[string]interface{}{
		"list":      treeList,
		"nameSpace": nameSpace,
	}
	return res
}

// LoadTableInfo 加载表信息
func (s *SysTableInfoService) LoadTableInfo(c *gin.Context) *response.WebResponseContent {
	var req request.LoadTableInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.Error("操作失败")
	}
	// tableId初始化
	if req.TableId <= 0 {
		req.TableId = InitTable(req.ParentId, req.TableName, req.ColumnCNName, req.NameSpace, req.FolderName, req.TableId, req.IsTreeLoad, req.DBServer)
		if req.TableId == -1 {
			return response.Error("获取表结构失败，请检查表名或数据库连接是否正确")
		}
	}
	var info system.SysTableInfo
	err := global.GVA_DB.
		Model(&system.SysTableInfo{}).
		Where("Table_Id = ?", req.TableId).
		Preload("TableColumns", func(db *gorm.DB) *gorm.DB {
			return db.Order("OrderNo ASC")
		}).
		First(&info).Error
	if err != nil {
		return response.Error(err.Error())
	}
	return response.Ok("", info)
}
func InitTable(
	parentId int,
	tableName, columnCNName, nameSpace, folderName string,
	tableId int,
	isTreeLoad bool,
	dbServer string,
) int {
	if isTreeLoad {
		return tableId
	}
	if tableName == "" {
		return -1
	}

	// 查询是否已存在
	var existInfo system.SysTableInfo
	err := global.GVA_DB.Where("table_name = ?", tableName).First(&existInfo).Error
	if err == nil && existInfo.TableId > 0 {
		return existInfo.TableId
	}

	// 构造表信息
	tableInfo := system.SysTableInfo{
		ParentId:     parentId,
		ColumnCNName: columnCNName,
		CnName:       columnCNName,
		Table_Name:   tableName,
		Namespace:    nameSpace,
		FolderName:   folderName,
		Enable:       1,
		DBServer:     dbServer,
	}

	// 查询表字典信息 没有就直接结束掉
	columns := GetTableColumns(tableName, dbServer)
	if len(columns) == 0 {
		return -1
	}

	// 设置顺序和编辑行号
	for i, col := range columns {
		col.OrderNo = i + 1
		col.EditRowNo = i + 1
	}

	SetMaxLength(columns)

	// 用事务插入主表和子表
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tableInfo).Error; err != nil {
			return err
		}
		for _, col := range columns {
			col.TableId = tableInfo.TableId
			if err := tx.Create(col).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return -1
	}
	return tableInfo.TableId
}
func SetMaxLength(columns []*system.SysTableColumn) {
	for _, x := range columns {
		switch {
		case x.ColumnName == "DateTime":
			x.ColumnWidth = 150
		case x.ColumnName == "Modifier" || x.ColumnName == "Creator":
			x.ColumnWidth = 130
		case x.ColumnName == "CreateID" || x.ColumnName == "ModifyID":
			x.ColumnWidth = 80
		case x.Maxlength > 200:
			x.ColumnWidth = 220
		case x.Maxlength > 110 && x.Maxlength <= 200:
			x.ColumnWidth = 180
		case x.Maxlength > 60 && x.Maxlength <= 110:
			x.ColumnWidth = 120
		default:
			x.ColumnWidth = 110
		}
	}
}

// GetTableColumns 获取表字段信息
func GetTableColumns(tableName, dbServer string) []*system.SysTableColumn {
	sql := GetCurrentSql(tableName, dbServer)
	if sql == "" {
		return []*system.SysTableColumn{}
	}
	var columns []*system.SysTableColumn
	if err := global.GetGlobalDBByDBName(dbServer).Raw(sql).Scan(&columns).Error; err != nil {
		return []*system.SysTableColumn{}
	}
	return columns
}

// GetCurrentSql 根据数据库类型返回查询表结构的SQL
func GetCurrentSql(tableName, dbService string) string {
	db := global.GetGlobalDBByDBName(dbService)
	switch db.Dialector.(type) {
	case *mysql.Dialector:
		return GetMySqlStructure(tableName, dbService)
	case *postgres.Dialector:
		return GetPgSqlStructure(tableName)
	case *sqlserver.Dialector:
		return GetSqlServerStructure(tableName)
	default:
		return ""
	}
}

func GetMySqlStructure(tableName, dbService string) string {
	return fmt.Sprintf(`SELECT DISTINCT
		Column_Name AS ColumnName,
		'%s' as tableName,
		Column_Comment AS ColumnCnName,
		CASE
			WHEN data_type IN('BIT', 'BOOL', 'bit', 'bool') THEN 'bool'
			WHEN data_type IN('smallint','SMALLINT') THEN 'int16'
			WHEN data_type IN('tinyint','TINYINT') THEN 'int8'
			WHEN data_type IN('MEDIUMINT','mediumint', 'int','INT','year', 'Year') THEN 'int32'
			WHEN data_type IN('BIGINT','bigint') THEN 'int64'
			WHEN data_type IN('FLOAT', 'DOUBLE', 'DECIMAL','float', 'double', 'decimal') THEN 'float64'
			WHEN data_type IN('CHAR', 'VARCHAR', 'TINY TEXT', 'TEXT', 'MEDIUMTEXT', 'LONGTEXT', 'TINYBLOB', 'BLOB', 'MEDIUMBLOB', 'LONGBLOB', 'Time','char', 'varchar', 'tiny text', 'text', 'mediumtext', 'longtext', 'tinyblob', 'blob', 'mediumblob', 'longblob', 'time') THEN 'string'
			WHEN data_type IN('Date', 'DateTime', 'TimeStamp','date', 'datetime', 'timestamp') THEN 'time.Time'
			ELSE 'string'
		END AS ColumnType,
		CASE WHEN CHARACTER_MAXIMUM_LENGTH>8000 THEN 0 ELSE CHARACTER_MAXIMUM_LENGTH END AS Maxlength,
		CASE WHEN COLUMN_KEY <> '' THEN 1 ELSE 0 END AS IsKey,
		CASE WHEN Column_Name IN('CreateID', 'ModifyID', '') OR COLUMN_KEY<> '' THEN 0 ELSE 1 END AS IsDisplay,
		1 AS IsColumnData,
		120 AS ColumnWidth,
		0 AS OrderNo,
		CASE WHEN IS_NULLABLE = 'NO' THEN 0 ELSE 1 END AS IsNull,
		CASE WHEN COLUMN_KEY <> '' THEN 1 ELSE 0 END AS IsReadDataset
	FROM information_schema.COLUMNS
	WHERE table_name = '%s' %s`,
		tableName,
		tableName,
		fmt.Sprintf(" and table_schema = '%s' ", dbService),
	)
}

func GetSqlServerStructure(tableName string) string {
	return fmt.Sprintf(`
	SELECT TableName,
		LTRIM(RTRIM(ColumnName)) AS ColumnName,
		ColumnCNName,
		CASE WHEN ColumnType = 'uniqueidentifier' THEN 'string' -- Go: uuid通常用string
			 WHEN ColumnType IN('smallint', 'INT') THEN 'int32' -- Go: int32
			 WHEN ColumnType = 'BIGINT' THEN 'int64' -- Go: int64
			 WHEN ColumnType IN('CHAR', 'VARCHAR', 'NVARCHAR', 'text', 'xml', 'varbinary', 'image') THEN 'string' -- Go: string
			 WHEN ColumnType IN('tinyint') THEN 'int8' -- Go: int8
			 WHEN ColumnType IN('bit') THEN 'bool' -- Go: bool
			 WHEN ColumnType IN('time', 'date', 'DATETIME', 'smallDATETIME') THEN 'time.Time' -- Go: time.Time
			 WHEN ColumnType IN('smallmoney', 'DECIMAL', 'numeric', 'money') THEN 'float64' -- Go: float64
			 WHEN ColumnType = 'float' THEN 'float64' -- Go: float64
			 ELSE 'string' -- Go: string
		END ColumnType,
		CASE WHEN ColumnType IN ('NVARCHAR','NCHAR') THEN [Maxlength]/2 ELSE [Maxlength] END [Maxlength],
		IsKey,
		CASE WHEN ColumnName IN('CreateID', 'ModifyID', '') OR IsKey = 1 THEN 0 ELSE 1 END AS IsDisplay,
		1 AS IsColumnData,
		CASE WHEN ColumnType IN('time', 'date', 'DATETIME', 'smallDATETIME') THEN 150
			 WHEN ColumnName IN('Modifier', 'Creator') THEN 130
			 WHEN ColumnType IN('int', 'bigint') OR ColumnName IN('CreateID', 'ModifyID', '') THEN 80
			 WHEN [Maxlength] < 110 AND [Maxlength] > 60 THEN 120
			 WHEN [Maxlength] < 200 AND [Maxlength] >= 110 THEN 180
			 WHEN [Maxlength] > 200 THEN 220
			 ELSE 110
		   END AS ColumnWidth,
		0 AS OrderNo,
		t.[IsNull] AS [IsNull],
		CASE WHEN IsKey = 1 THEN 1 ELSE 0 END IsReadDataset,
		CASE WHEN IsKey!=1 AND t.[IsNull] = 0 THEN 0 ELSE NULL END AS EditColNo
	FROM (
		SELECT obj.name AS TableName,
			col.name AS ColumnName,
			CONVERT(NVARCHAR(100),ISNULL(ep.[value], '')) AS ColumnCNName,
			t.name AS ColumnType,
			CASE WHEN col.length<1 THEN 0 ELSE col.length END AS [Maxlength],
			CASE WHEN EXISTS (
				SELECT 1
				FROM dbo.sysindexes si
				INNER JOIN dbo.sysindexkeys sik ON si.id = sik.id AND si.indid = sik.indid
				INNER JOIN dbo.syscolumns sc ON sc.id = sik.id AND sc.colid = sik.colid
				INNER JOIN dbo.sysobjects so ON so.name = si.name AND so.xtype = 'PK'
				WHERE sc.id = col.id AND sc.colid = col.colid
			) THEN 1 ELSE 0 END AS IsKey,
			CASE WHEN col.isnullable = 1 THEN 1 ELSE 0 END AS [IsNull],
			col.colorder
		FROM dbo.syscolumns col
		LEFT JOIN dbo.systypes t ON col.xtype = t.xusertype
		INNER JOIN dbo.sysobjects obj ON col.id = obj.id AND obj.xtype IN ('U','V')
		LEFT JOIN dbo.syscomments comm ON col.cdefault = comm.id
		LEFT JOIN sys.extended_properties ep ON col.id = ep.major_id AND col.colid = ep.minor_id AND ep.name = 'MS_Description'
		LEFT JOIN sys.extended_properties epTwo ON obj.id = epTwo.major_id AND epTwo.minor_id = 0 AND epTwo.name = 'MS_Description'
		WHERE obj.name = '%s'
	) AS t
	ORDER BY t.colorder
	`, tableName)
}

func GetPgSqlStructure(tableName string) string {
	return fmt.Sprintf(`
SELECT 
	MM."TableName", 
	MM."ColumnName", 
	MM."ColumnCNName", 
	CASE 
		WHEN MM."ColumnType" = 'uuid' THEN 'string' -- Go: uuid通常用string
		WHEN MM."ColumnType" = 'short' THEN 'int16'
		WHEN MM."ColumnType" = 'int' THEN 'int32'
		WHEN MM."ColumnType" = 'long' THEN 'int64'
		WHEN MM."ColumnType" = 'string' THEN 'string'
		WHEN MM."ColumnType" = 'bool' THEN 'bool'
		WHEN MM."ColumnType" = 'DateTime' THEN 'time.Time'
		WHEN MM."ColumnType" = 'decimal' THEN 'float64'
		WHEN MM."ColumnType" = 'float' THEN 'float64'
		ELSE 'string'
	END AS "ColumnType",
	MM."Maxlength", 
	MM."IsKey", 
	MM."IsDisplay", 
	MM."IsColumnData", 
	CASE 
		WHEN MM."ColumnType" = 'time.Time' THEN 150  
		WHEN MM."ColumnType" = 'int32' THEN 80  
		WHEN MM."Maxlength" < 110 AND MM."Maxlength" > 60 THEN 120  
		WHEN MM."Maxlength" < 200 AND MM."Maxlength" >= 110 THEN 180  
		WHEN MM."Maxlength" > 200 THEN 220 ELSE 110  
	END AS "ColumnWidth", 
	MM."OrderNo", 
	CASE WHEN MM."IsKey"=1 or lower(MM."IsNull")='no' then 0 else 1 end as "IsNull", 
	MM."IsReadDataset", 
	MM."EditColNo"  
FROM (
	SELECT 
		col.TABLE_NAME AS "TableName", 
		col.COLUMN_NAME AS "ColumnName", 
		attr.description AS "ColumnCNName", 
		CASE 
			WHEN col.udt_name = 'uuid' THEN 'string'  
			WHEN col.udt_name IN ('int2') THEN 'int16'  
			WHEN col.udt_name IN ('int4') THEN 'int32'  
			WHEN col.udt_name = 'int8' THEN 'int64'  
			WHEN col.udt_name = 'BIGINT' THEN 'int64'  
			WHEN col.udt_name IN ('char', 'varchar', 'text', 'xml', 'bytea') THEN 'string'  
			WHEN col.udt_name IN ('bool') THEN 'bool'  
			WHEN col.udt_name IN ('date','timestamp') THEN 'time.Time'  
			WHEN col.udt_name IN ('decimal', 'money','numeric') THEN 'float64'  
			WHEN col.udt_name IN ('float4', 'float8') THEN 'float64' ELSE 'string'  
		END "ColumnType", 
		CASE 
			WHEN col.udt_name = 'varchar' THEN col.character_maximum_length  
			WHEN col.udt_name IN ('int2', 'int4', 'int8', 'float4', 'float8') THEN col.numeric_precision ELSE 1024  
		END "Maxlength", 
		CASE WHEN keyTable.IsKey = 1 THEN 1 ELSE 0 END "IsKey", 
		CASE WHEN keyTable.IsKey = 1 THEN 0 ELSE 1 END "IsDisplay", 
		1 AS "IsColumnData", 
		0 AS "OrderNo", 
		col.is_nullable AS "IsNull", 
		CASE WHEN keyTable.IsKey = 1 THEN 1 ELSE 0 END "IsReadDataset", 
		CASE WHEN keyTable.IsKey IS NULL AND col.is_nullable = 'NO' THEN 0 ELSE NULL END AS "EditColNo"  
	FROM information_schema.COLUMNS col  
	LEFT JOIN (
		SELECT col_description(a.attrelid,a.attnum) as description,a.attname as name 
		FROM pg_class as c,pg_attribute as a  
		WHERE lower(c.relname) = lower('%s') and a.attrelid = c.oid and a.attnum>0 
	) as attr on attr.name=col.COLUMN_NAME 
	LEFT JOIN (
		SELECT 
			pg_attribute.attname AS colname, 
			1 AS IsKey  
		FROM pg_constraint 
		INNER JOIN pg_class ON pg_constraint.conrelid = pg_class.oid 
		INNER JOIN pg_attribute ON pg_attribute.attrelid = pg_class.oid  
			AND pg_attribute.attnum = pg_constraint.conkey [1]  
		WHERE lower(pg_class.relname) = lower('%s')  
			AND pg_constraint.contype = 'p'  
	) keyTable ON col.COLUMN_NAME = keyTable.colname  
	WHERE lower(TABLE_NAME) = lower('%s')  
	ORDER BY ordinal_position  
) MM; 
	`, tableName, tableName, tableName)
}

// GetConnectionString 获取连接字符串 根据 gorm.DB 获取 DSN
func GetConnectionString(dbConnection string) string {
	db := global.GetGlobalDBByDBName(dbConnection)
	switch dial := db.Dialector.(type) {
	case *mysql.Dialector:
		return dial.DSN
	case *postgres.Dialector:
		return dial.DSN
	case *sqlserver.Dialector:
		return dial.DSN
	default:
		return ""
	}
}

type Field struct {
	Name         string // Go struct 字段名
	Type         string // Go 字段类型
	GormTag      string // gorm标签
	JsonName     string // json标签
	Nullable     bool   // 是否可为空
	Editable     bool   // 是否可编辑
	Display      bool   // 是否显示
	Key          bool   // 是否主键
	ColumnCNName string // 中文名
	ColumnName   string // 原始字段名
}

type TemplateData struct {
	PackageName string
	StructName  string
	TableName   string
	CnName      string
	Fields      []Field
}

// CreateEntityModel 生成model文件
func (s *SysTableInfoService) CreateEntityModel(req system.SysTableInfo) (TemplateData, error) {
	tableId := req.TableId
	// 获取表信息
	var tableInfo system.SysTableInfo
	err := global.GVA_DB.
		Model(&system.SysTableInfo{}).
		Where("Table_Id = ?", tableId).
		Preload("TableColumns", func(db *gorm.DB) *gorm.DB {
			return db.Order("OrderNo ASC")
		}).
		First(&tableInfo).Error
	if err != nil {
		return TemplateData{}, err
	}

	fields := make([]Field, 0, len(tableInfo.TableColumns))
	for _, col := range tableInfo.TableColumns {
		goFieldName := utils.CamelCase(col.ColumnName)
		goType := utils.GoTypeWithNull(col.ColumnType, col.IsNull)
		meta := generateColumnMeta(col)
		fields = append(fields, Field{
			Name:         goFieldName,
			Type:         goType,
			GormTag:      fmt.Sprintf("column:%s", goFieldName),
			JsonName:     goFieldName,
			Nullable:     meta.Nullable,
			Editable:     meta.Editable,
			Display:      meta.Display,
			Key:          col.IsKey == 1,
			ColumnCNName: col.ColumnCNName,
			ColumnName:   col.ColumnName,
		})
	}

	data := TemplateData{
		PackageName: "model",
		StructName:  utils.CamelCase(tableInfo.Table_Name),
		TableName:   tableInfo.Table_Name,
		CnName:      tableInfo.CnName,
		Fields:      fields,
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		return TemplateData{}, err
	}
	dirPath := filepath.Join(projectRoot, "model", tableInfo.FolderName)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return TemplateData{}, err
	}
	filePath := filepath.Join(dirPath, "test_template.go")
	f, err := os.Create(filePath)
	if err != nil {
		return TemplateData{}, err
	}
	defer func() {
		cerr := f.Close()
		if cerr != nil {
			fmt.Fprintf(os.Stderr, "close file error: %v\n", cerr)
		}
	}()

	tmplPath := filepath.Join(projectRoot, "tmpl", "model.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return TemplateData{}, err
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return TemplateData{}, err
	}

	// 自动 gofmt 格式化
	exec.Command("gofmt", "-w", filePath).Run()

	return data, nil
}

// columnMeta 字段元属性
type columnMeta struct {
	Nullable bool
	Editable bool
	Display  bool
	Sortable bool
}

// generateColumnMeta 根据 SysTableColumn 自动生成字段元属性
func generateColumnMeta(col system.SysTableColumn) columnMeta {
	return columnMeta{
		Nullable: col.IsNull == 1,
		Editable: col.IsKey == 0 && col.ApiInPut == 1,
		Display:  col.IsDisplay == 1,
		Sortable: col.Sortable == 1,
	}
}
