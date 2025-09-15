package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
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

// CreateModel 生成model文件
// @Tags     SysTableInfo
// @Summary  生成model文件
// @Produce  application/json
// @Param    data  body      system.SysTableInfo  true "参数"
// @Success 200 {object} response.Response{data=response.WebResponseContent} "生成model文件"
// @Router   /builder/createModel [post]
func (s *SysTableInfoApi) CreateModel(c *gin.Context) {
	var req system.SysTableInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WebResponse(response.Error("参数错误: tableId 必填"), c)
		return
	}
	res, err := Service().CreateModel(req)
	if err != nil {
		response.WebResponse(response.Error(err.Error()), c)
		return
	}
	c.JSON(http.StatusOK, res)
}

// CreateServices 生成service文件
// @Tags     SysTableInfo
// @Summary  生成service文件
// @Produce  application/json
// @Param    data  body      system.SysTableInfo  true "参数"
// @Success 200 {object} response.Response{data=response.WebResponseContent} "生成service文件"
// @Router   /builder/createServices [post]
func (s *SysTableInfoApi) CreateServices(c *gin.Context) {
	var req system.SysTableInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WebResponse(response.Error("参数错误: tableId 必填"), c)
		return
	}
	res, err := Service().CreateServices(req)
	if err != nil {
		response.WebResponse(response.Error(err.Error()), c)
		return
	}
	c.JSON(http.StatusOK, res)
}
