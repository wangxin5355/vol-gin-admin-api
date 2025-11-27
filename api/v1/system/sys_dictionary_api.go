package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/service/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

type SysDictionaryApi struct{}

func (api *SysDictionaryApi) SysDictionaryService() *system.DictionaryService {
	return service.ServiceInstances.DictionaryService
}

// GetPageData
// @Tags     SysDictionaryApi
// @Summary  获取分页数据
// @Produce  application/json
// @Param    options  body	  request.PageDataOptions  true  "分页数据选项"
// @Success 200 {object} response.Response{data=[]system.SysDictionary} "返回分页数据"
// @Router   /api/Sys_Dictionary/GetPageData [post]
func (api *SysDictionaryApi) GetPageData(c *gin.Context) {
	param, err := utils.BindJsonToPageDataOptions(c)
	if err != nil {
		return
	}
	data := api.SysDictionaryService().GetPageData(param)
	response.OkWithPageData(data.Rows, data.Summary, data.Total, c)
}

func (api *SysDictionaryApi) GetDetailPage(c *gin.Context) {
	param, err := utils.BindJsonToPageDataOptions(c)
	if err != nil {
		return
	}
	data := api.SysDictionaryService().GetDetailPage(param)
	response.OkWithPageData(data.Rows, data.Summary, data.Total, c)
}

// GetBuilderDictionary
// @Tags     SysDictionaryApi
// @Summary   代码生成器获取所有字典项(超级管理权限)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success  200   {object}  []string  "返回所有字典"
// @Router   /api/Sys_Dictionary/GetBuilderDictionary [get]
func (api *SysDictionaryApi) GetBuilderDictionary(c *gin.Context) {
	dicNos := service.ServiceInstances.DictionaryService.GetBuilderDictionary()
	c.JSON(http.StatusOK, dicNos)
}

// GetVueDictionary
// @Tags     SysDictionaryApi
// @Summary   获取vue页面需要的字典
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     dicNos  body    []string  true  "字典编号数组"
// @Success  200   {object}  []map[string]interface{}  "返回vue字典数据"
// @Router   /api/Sys_Dictionary/GetVueDictionary [post]
func (api *SysDictionaryApi) GetVueDictionary(c *gin.Context) {
	var dicNos []string
	if err := c.ShouldBindJSON(&dicNos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vueDict := service.ServiceInstances.DictionaryService.GetVueDictionar(dicNos)
	c.JSON(http.StatusOK, vueDict)
}

// Add
// @Tags     SysDictionaryApi
// @Summary  新增字典数据
// @Produce  application/json
// @Param    sysDictionary  body	  system.SysDictionary  true  "字典数据"
// @Success 200 {object} response.Response{data=system.SysDictionary} "返回新增的字典数据"
// @Router   /api/Sys_Dictionary/Add [post]
func (api *SysDictionaryApi) Add(c *gin.Context) {
	param, err := utils.BindJsonToSaveModel(c)
	if err != nil {
		return
	}
	data := api.SysDictionaryService().Add(c, param)
	//response.OkWithData(data, c)
	c.JSON(http.StatusOK, data)
}

// Update
// @Tags     SysDictionaryApi
// @Summary  更新字典数据
// @Produce  application/json
// @Param    sysDictionary  body	  system.SysDictionary  true  "字典数据"
// @Success 200 {object} response.Response{data=system.SysDictionary} "返回更新后的字典数据"
// @Router   /api/Sys_Dictionary/update [post]
func (api *SysDictionaryApi) Update(c *gin.Context) {
	param, err := utils.BindJsonToSaveModel(c)
	if err != nil {
		return
	}
	data := api.SysDictionaryService().Update(c, param)
	//response.OkWithData(data, c)
	c.JSON(http.StatusOK, data)
}

func (api *SysDictionaryApi) Del(c *gin.Context) {
	var keys []any
	err := c.ShouldBindJSON(&keys)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data := api.SysDictionaryService().Del(c, keys)
	c.JSON(http.StatusOK, data)
}
