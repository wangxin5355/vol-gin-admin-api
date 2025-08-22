package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/core/base/provider"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

type TestApi struct {
}

var testService = service.ServiceGroupApp.TestServiceGroup.TestService

// GetPageData
// @Tags     TestApi
// @Summary  获取分页数据
// @Produce  application/json
// @Param    options  body	  provider.PageDataOptions  true  "分页数据选项"
// @Success  200   {object}  response.Response{data=provider.PageGridData[system.SysUser],msg=string}  "返回分页数据"
// @Router   /test/GetPageData [post]
func (b *TestApi) GetPageData(c *gin.Context) {
	param, ok := utils.BindJSONParam[provider.PageDataOptions](c)
	if !ok {
		return
	}
	data := testService.GetPageData(param)
	response.OkWithData(data, c)
}
