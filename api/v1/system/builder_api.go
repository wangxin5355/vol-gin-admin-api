package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/service"
)

type BuilderApi struct{}

// GetTableTree
// @Tags     BuilderApi
// @Summary  获取代码生成树
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success  200   {object}  response.TableTreeResp  "返回所有的设计表树"
// @Router   /builder/GetTableTree [get]
func (api *BuilderApi) GetTableTree(c *gin.Context) {
	//Sys_TableInfoService获取表信息 ，golang没有类库不需要返回，路径想其他方式
	res := service.ServiceInstances.SysTableInfoService.GetTableTree()
	c.JSON(http.StatusOK, res)
}

// CreateModel
// @Tags     BuilderApi
// @Summary  创建Model文件
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param    data  body      map[string]interface{}  true  "参数"
// @Success  200   {string}  string                  "Model创建成功"
// @Router   /builder/CreateModel [post]
func (api *BuilderApi) CreateModel(c *gin.Context) {
	_, err := service.ServiceInstances.SysTableInfoService.CreateModel(c)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Model创建成功")
}

// CreateServices 生成service文件
// @Tags     SysTableInfo
// @Summary  生成service文件
// @Produce  application/json
// @Param    data  body      system.SysTableInfo  true "参数"
// @Success 200 {object} response.Response{data=response.WebResponseContent} "生成service文件"
// @Router   /builder/createServices [post]
func (api *BuilderApi) CreateServices(c *gin.Context) {
	var req system.SysTableInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WebResponse(response.Error("参数错误: tableId 必填"), c)
		return
	}
	_, err := Service().CreateServices(req)
	if err != nil {
		response.WebResponse(response.Error(err.Error()), c)
		return
	}
	c.JSON(http.StatusOK, "Service创建成功")
}
