package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/core/base/provider"
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
func BindJsonToPageDataOptions(c *gin.Context) (provider.PageDataOptions, error) {
	return BindJSONParam[provider.PageDataOptions](c)
}
