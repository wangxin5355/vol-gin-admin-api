package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/service/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

type SysLogApi struct{}

func (api *SysLogApi) SysLogService() *system.SysLogService {
	return service.ServiceInstances.SysLogService
}

// GetPageData
// @Tags     SysLogApi
// @Summary  获取分页数据
// @Produce  application/json
// @Param    options  body	  request.PageDataOptions  true  "分页数据选项"
// @Success 200 {object} response.Response{data=[]system.Sys_Log} "返回分页数据"
// @Router   /test/GetPageData [post]
func (api *SysLogApi) GetPageData(c *gin.Context) {
	param, err := utils.BindJsonToPageDataOptions(c)
	if err != nil {
		return
	}
	data := api.SysLogService().GetPageData(param)
	response.OkWithPageData(data.Rows, data.Summary, data.Total, c)
}

// Add
// @Tags     SysLogApi
// @Summary  添加数据
// @Produce  application/json
// @Param    saveModel  body	  request.SaveModel  true  "添加数据"
// @Success 200 {object} response.Response{data=string} "添加数据"
// @Router   /test/Add [post]
func (api *SysLogApi) Add(c *gin.Context) {
	param, err := utils.BindJsonToSaveModel(c)
	if err != nil {
		return
	}
	data := api.SysLogService().Add(c, param)
	response.OkWithData(data, c)
}

// Update
// @Tags     SysLogApi
// @Summary  更新数据
// @Produce  application/json
// @Param    saveModel  body	  request.SaveModel  true  "更新数据"
// @Success 200 {object} response.Response{data=string} "更新数据"
// @Router   /test/Update [post]
func (api *SysLogApi) Update(c *gin.Context) {
	param, err := utils.BindJsonToSaveModel(c)
	if err != nil {
		return
	}
	data := api.SysLogService().Update(c, param)
	response.OkWithData(data, c)

}

// Del
// @Tags     SysLogApi
// @Summary  删除数据
// @Produce  application/json
// @Param    keys  body	  []any  true  "删除数据"
// @Success 200 {object} response.Response{data=string} "删除数据"
// @Router   /test/Del [post]
func (api *SysLogApi) Del(c *gin.Context) {
	var keys []any
	err := c.ShouldBindJSON(&keys)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data := api.SysLogService().Del(c, keys)
	response.OkWithData(data, c)
}
