package code

import (
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
)

var SysTableInfoGroup struct {
	SysTableInfoService *SysTableInfoService
}

type SysTableInfoService struct {
}

// GetTableTree 获取表结构树形数据
func (s *SysTableInfoService) GetTableTree() *response.WebResponseContent {
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
		return response.Error("获取数据失败")
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
	return response.Ok("", treeList)
}
