package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/service/code"
)

type SysTableInfoApi struct {
}

func Service() *code.SysTableInfoService {
	return code.SysTableInfoGroup.SysTableInfoService
}

// GetTableTree 获取表结构树形数据
// @Tags     SysTableInfo
// @Summary  获取表结构树形数据
// @Produce  application/json
// @Success 200 {object} response.Response{data=response.WebResponseContent} "获取表结构树形数据"
// @Router   /sysTableInfo/getTableTree [get]
func (s *SysTableInfoApi) GetTableTree(c *gin.Context) {
	data := Service().GetTableTree()
	response.WebResponse(data, c)
}
