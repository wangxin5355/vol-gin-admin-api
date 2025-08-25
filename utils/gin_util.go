package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
)

// BindJSONParam 绑定传参到结构体
func BindJSONParam[T any](c *gin.Context) (T, error) {
	var param T
	if err := c.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("参数错误", c)
		return param, err
	}
	return param, nil
}

// BindJsonToPageDataOptions 将JSON绑定到分页数据选项
func BindJsonToPageDataOptions(c *gin.Context) (request.PageDataOptions, error) {
	return BindJSONParam[request.PageDataOptions](c)
}

// BindJsonToSaveModel 实体绑定
func BindJsonToSaveModel(c *gin.Context) (request.SaveModel, error) {
	return BindJSONParam[request.SaveModel](c)
}
