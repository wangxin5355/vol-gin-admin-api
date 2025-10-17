package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/system/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"net/http"
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
	list, nameSpace := service.ServiceInstances.TableInfoService.GetTableTree()
	c.JSON(http.StatusOK, response.TableTreeResp{
		List:      list,
		NameSpace: nameSpace,
	})
}
