package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

type TestApi struct {
}

func getServiceGroup() *service.TestServiceGroup {
	return &service.ServiceGroupApp.TestServiceGroup
}

// GetPageData
// @Tags     TestApi
// @Summary  获取分页数据
// @Produce  application/json
// @Param    options  body	  request.PageDataOptions  true  "分页数据选项"
// @Success 200 {object} response.Response{data=[]system.SysUser} "返回分页数据"
// @Router   /test/GetPageData [post]
func (b *TestApi) GetPageData(c *gin.Context) {
	param, err := utils.BindJsonToPageDataOptions(c)
	if err != nil {
		return
	}
	data := getServiceGroup().GetPageData(param)
	response.OkWithData(data, c)
}

// Add
// @Tags     TestApi
// @Summary  添加数据
// @Produce  application/json
// @Param    saveModel  body	  request.SaveModel  true  "添加数据"
// @Success 200 {object} response.Response{data=string} "添加数据"
// @Router   /test/Add [post]
func (b *TestApi) Add(c *gin.Context) {
	param, err := utils.BindJsonToSaveModel(c)
	if err != nil {
		return
	}
	data := getServiceGroup().Add(c, param)
	response.WebResponse(data, c)
}

// Update
// @Tags     TestApi
// @Summary  更新数据
// @Produce  application/json
// @Param    saveModel  body	  request.SaveModel  true  "更新数据"
// @Success 200 {object} response.Response{data=string} "更新数据"
// @Router   /test/Update [post]
func (b *TestApi) Update(c *gin.Context) {
	param, err := utils.BindJsonToSaveModel(c)
	if err != nil {
		return
	}
	data := getServiceGroup().Update(c, param)
	response.WebResponse(data, c)
}

// Del
// @Tags     TestApi
// @Summary  删除数据
// @Produce  application/json
// @Param    keys  body	  []any  true  "删除数据"
// @Success 200 {object} response.Response{data=string} "删除数据"
// @Router   /test/Del [post]
func (b *TestApi) Del(c *gin.Context) {
	var keys []any
	err := c.ShouldBindJSON(&keys)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data := getServiceGroup().Del(c, keys)
	response.WebResponse(data, c)
}

// GetCurrentUserInfo
// @Tags     TestApi
// @Summary  获取当前用户信息
// @Produce  application/json
// @Success 200 {object} response.Response{data=string} "获取当前用户
// @Router   /test/GetCurrentUserInfo [get]
func (b *TestApi) GetCurrentUserInfo(c *gin.Context) {
	data := getServiceGroup().GetCurrentUserInfo(c)
	response.WebResponse(data, c)
}
