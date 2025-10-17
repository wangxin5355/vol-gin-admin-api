package system

import (
	"encoding/json"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/dto"
	"log"
)

type TableInfoService struct {
}

// 获取生成配置的树开菜单
func (tableInfoService *TableInfoService) GetTableTree() (string, string) {
	var tableInfos []dto.TableInfo
	result := global.GVA_DB.Raw("SELECT Table_Id as Id,ParentId as PId,ParentId ,ColumnCNName as Name,OrderNo  FROM `sys_tableinfo` ORDER BY OrderNo").Scan(&tableInfos)
	if result.Error != nil {
		log.Fatal(result.Error)
		return "[]", "''"
	}
	var tableTreeList = getTableTreeList(tableInfos)
	data, err := json.Marshal(tableTreeList)
	if err != nil {
		panic(err)
	}
	return string(data), "" //go版本不返回命名空间。
}

// 获取tableTreeListdatas
func getTableTreeList(tables []dto.TableInfo) []dto.TableTreeListData {
	pids := make(map[int]struct{})
	for _, table := range tables {
		pids[table.PId] = struct{}{}
	}
	tableTreeListdatas := make([]dto.TableTreeListData, 0)
	for _, table := range tables {
		isParent := false
		if _, exists := pids[table.Id]; exists {
			isParent = true
		}
		tableTreeListdatas = append(tableTreeListdatas, dto.TableTreeListData{
			Id:       table.Id,
			PId:      table.PId,
			ParentId: table.ParentId,
			Name:     table.Name,
			IsParent: isParent,
		})
	}
	return tableTreeListdatas
}
