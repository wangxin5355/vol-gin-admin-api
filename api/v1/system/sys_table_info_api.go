package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Router   /builder/getTableTree [get]
func (s *SysTableInfoApi) GetTableTree(c *gin.Context) {
	data := Service().GetTableTree()
	c.JSON(http.StatusOK, data)
}

// LoadTableInfo 加载表信息
// @Tags     SysTableInfo
// @Summary  加载表信息
// @Produce  application/json
// @Param    data  body      request.LoadTableInfoReq  true  "加载表信息"
// @Success 200 {object} response.Response{data=response.WebResponseContent} "加载
// @Router   /builder/loadTableInfo [post]
func (s *SysTableInfoApi) LoadTableInfo(c *gin.Context) {
	res := Service().LoadTableInfo(c)
	c.JSON(http.StatusOK, res)
}
