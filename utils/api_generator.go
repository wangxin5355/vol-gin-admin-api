// utils/controller_generator.go
package utils

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type apiConfig struct {
	ModelName string
	BasePath  string
}

// GenerateApiFile 生成控制器文件
func GenerateApiFile(config apiConfig, outputPath string) error {
	tmpl := `package api

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "your-app/models"
    "your-app/services"
    "your-app/utils"
)

var (
    {{.ModelName | ToLower}}Service *services.BaseService
)

// Init{{.ModelName}}Controller 初始化{{.ModelName}}控制器
func Init{{.ModelName}}Controller(db *gorm.DB) {
    {{.ModelName | ToLower}}Service = services.NewBaseService(db)
}

// Create{{.ModelName}} 创建{{.ModelName}}
// @Summary 创建{{.ModelName}}
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param data body models.{{.ModelName}} true "{{.ModelName}}数据"
// @Success 200 {object} utils.Response
// @Router {{.BasePath}} [post]
func Create{{.ModelName}}(c *gin.Context) {
    var {{.ModelName | ToLower}} models.{{.ModelName}}
    if err := c.ShouldBindJSON(&{{.ModelName | ToLower}}); err != nil {
        utils.ResponseError(c, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }
    
    if err := {{.ModelName | ToLower}}Service.Create(&{{.ModelName | ToLower}}); err != nil {
        utils.ResponseError(c, http.StatusInternalServerError, "创建失败: "+err.Error())
        return
    }
    
    utils.ResponseSuccess(c, {{.ModelName | ToLower}})
}

// Get{{.ModelName}}List 获取{{.ModelName}}列表
// @Summary 获取{{.ModelName}}列表
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "关键词"
// @Success 200 {object} utils.Response
// @Router {{.BasePath}} [get]
func Get{{.ModelName}}List(c *gin.Context) {
    var pageInfo models.PageInfo
    if err := c.ShouldBindQuery(&pageInfo); err != nil {
        utils.ResponseError(c, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }
    
    // 设置默认值
    if pageInfo.Page <= 0 {
        pageInfo.Page = 1
    }
    if pageInfo.PageSize <= 0 {
        pageInfo.PageSize = 10
    }
    
    var {{.ModelName | ToLower}}List []models.{{.ModelName}}
    result, err := {{.ModelName | ToLower}}Service.GetList(&{{.ModelName | ToLower}}List, pageInfo)
    if err != nil {
        utils.ResponseError(c, http.StatusInternalServerError, "查询失败: "+err.Error())
        return
    }
    
    utils.ResponseSuccess(c, result)
}

// Get{{.ModelName}}ByID 根据ID获取{{.ModelName}}
// @Summary 根据ID获取{{.ModelName}}
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param id path int true "{{.ModelName}} ID"
// @Success 200 {object} utils.Response
// @Router {{.BasePath}}/{id} [get]
func Get{{.ModelName}}ByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ResponseError(c, http.StatusBadRequest, "ID参数错误")
        return
    }
    
    var {{.ModelName | ToLower}} models.{{.ModelName}}
    if err := {{.ModelName | ToLower}}Service.GetByID(&{{.ModelName | ToLower}}, uint(id)); err != nil {
        utils.ResponseError(c, http.StatusNotFound, "记录不存在")
        return
    }
    
    utils.ResponseSuccess(c, {{.ModelName | ToLower}})
}

// Update{{.ModelName}} 更新{{.ModelName}}
// @Summary 更新{{.ModelName}}
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param id path int true "{{.ModelName}} ID"
// @Param data body models.{{.ModelName}} true "{{.ModelName}}数据"
// @Success 200 {object} utils.Response
// @Router {{.BasePath}}/{id} [put]
func Update{{.ModelName}}(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ResponseError(c, http.StatusBadRequest, "ID参数错误")
        return
    }
    
    var {{.ModelName | ToLower}} models.{{.ModelName}}
    if err := {{.ModelName | ToLower}}Service.GetByID(&{{.ModelName | ToLower}}, uint(id)); err != nil {
        utils.ResponseError(c, http.StatusNotFound, "记录不存在")
        return
    }
    
    if err := c.ShouldBindJSON(&{{.ModelName | ToLower}}); err != nil {
        utils.ResponseError(c, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }
    
    if err := {{.ModelName | ToLower}}Service.Update(&{{.ModelName | ToLower}}); err != nil {
        utils.ResponseError(c, http.StatusInternalServerError, "更新失败: "+err.Error())
        return
    }
    
    utils.ResponseSuccess(c, {{.ModelName | ToLower}})
}

// Delete{{.ModelName}} 删除{{.ModelName}}
// @Summary 删除{{.ModelName}}
// @Tags {{.ModelName}}
// @Accept json
// @Produce json
// @Param id path int true "{{.ModelName}} ID"
// @Success 200 {object} utils.Response
// @Router {{.BasePath}}/{id} [delete]
func Delete{{.ModelName}}(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ResponseError(c, http.StatusBadRequest, "ID参数错误")
        return
    }
    
    var {{.ModelName | ToLower}} models.{{.ModelName}}
    if err := {{.ModelName | ToLower}}Service.Delete(&{{.ModelName | ToLower}}, uint(id)); err != nil {
        utils.ResponseError(c, http.StatusInternalServerError, "删除失败: "+err.Error())
        return
    }
    
    utils.ResponseSuccess(c, nil)
}

{{range .CustomHandlers}}
// {{.Name}} {{.Description}}
// @Summary {{.Description}}
// @Tags {{$.ModelName}}
// @Accept json
// @Produce json
{{range .Params}}
// @Param {{.}} 
{{end}}
// @Success 200 {object} utils.Response
// @Router {{$.BasePath}}{{.Path}} [{{.Method}}]
func {{.Name}}(c *gin.Context) {
    // TODO: 实现 {{.Description}} 逻辑
    utils.ResponseSuccess(c, gin.H{
        "message": "{{.Description}} - 待实现",
    })
}
{{end}}
`

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t, err := template.New("controller").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 创建目录
	dir := outputPath[:strings.LastIndex(outputPath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	return t.Execute(file, config)
}
