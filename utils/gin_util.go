package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
)

// 绑定传参到结构体
func BindJSONParam[T any](c *gin.Context) (T, bool) {
	var param T
	if err := c.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("参数错误", c)
		return param, false
	}
	return param, true
}
