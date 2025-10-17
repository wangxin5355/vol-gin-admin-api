package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/service/test"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

type TestApi struct {
	//base.BaseApi[system.SysUser, base.ServiceInterface[system.SysUser]]
	//base.BaseApi[partial.TestTemplateEntity, base.ServiceInterface[partial.TestTemplateEntity]]
	//base.BaseApi[partial.TestTemplateEntity, *test.TestService]
}

// Service 重写 Service 方法
func (b *TestApi) Service() *test.TestService {
	return service.ServiceInstances.TestService
}

//// GetPageData
//// @Tags     TestApi
//// @Summary  获取分页数据
//// @Produce  application/json
//// @Param    options  body	  request.PageDataOptions  true  "分页数据选项"
//// @Success 200 {object} response.Response{data=[]system.SysUser} "返回分页数据"
//// @Router   /test/GetPageData [post]
//func (b *TestApi) GetPageData(c *gin.Context) {
//	param, err := utils.BindJsonToPageDataOptions(c)
//	if err != nil {
//		return
//	}
//	data := Service().GetPageData(param)
//	response.OkWithData(data, c)
//}

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
	data := b.Service().Add(c, param)
	response.OkWithData(data, c)
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
	data := b.Service().Update(c, param)
	response.OkWithData(data, c)

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
	data := b.Service().Del(c, keys)
	response.OkWithData(data, c)

}

// GetCurrentUserInfo
// @Tags     TestApi
// @Summary  获取当前用户信息
// @Produce  application/json
// @Success 200 {object} response.Response{data=string} "获取当前用户
// @Router   /test/GetCurrentUserInfo [get]
func (b *TestApi) GetCurrentUserInfo(c *gin.Context) {
	data := b.Service().GetCurrentUserInfo(c)
	response.WebResponse(data, c)
}

// RedisTest
// @Tags     TestApi
// @Summary  测试 Redis
// @Produce  application/json
// @Success 200 {object} response.Response{data=string} "测试 Redis
// @Router   /test/RedisTest [get]
func (b *TestApi) RedisTest(c *gin.Context) {
	data := b.Service().RedisTest()
	response.WebResponse(data, c)
}
