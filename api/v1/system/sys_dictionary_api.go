package system

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"net/http"
)

type SysDictionaryApi struct{}

// GetBuilderDictionary
// @Tags     SysDictionaryApi
// @Summary   代码生成器获取所有字典项(超级管理权限)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success  200   {object}  []string  "返回所有字典"
// @Router   /Sys_Dictionary/GetBuilderDictionary [get]
func (api *SysDictionaryApi) GetBuilderDictionary(c *gin.Context) {
	dicNos := service.ServiceInstances.DictionaryService.GetBuilderDictionary()
	c.JSON(http.StatusOK, dicNos)
}
