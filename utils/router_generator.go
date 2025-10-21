// utils/router_generator.go
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// RouteConfig 路由配置
type RouteConfig struct {
	ModelName string // 模型名
	BasePath  string // 基础路径
}

// GenerateRouteFile 生成路由注册文件
func GenerateRouteFile(config RouteConfig, outputPath string) error {
	tmpl := `package router

import (
    "github.com/gin-gonic/gin"
api "github.com/wangxin5355/vol-gin-admin-api/api/v1"
    "github.com/wangxin5355/vol-gin-admin-api/middleware"
)

// Register{{.ModelName}}Routes 注册{{.ModelName}}路由  通过反射，找到
func Register{{.ModelName}}Routes(router *gin.RouterGroup) {
    {{$modelName := .ModelName}}
    // {{.ModelName}}路由组 所有api授权 ，如果不需要授权的，请单独传入不授权的RouterGroup
    {{.ModelName | ToLower}}Group := router.Group("{{.BasePath}}")

    {{.ModelName | ToLower}}Group.Use(middleware.JWTAuth())
    {
        // 自动生成的CRUD路由
        {{.ModelName | ToLower}}Group.POST("Create{{.ModelName}}", {{.ModelName | ToLower}}Api.Create{{.ModelName}})
        {{.ModelName | ToLower}}Group.GET("Get{{.ModelName}}List", {{.ModelName | ToLower}}Api.Get{{.ModelName}}List)
    }
}

// Get{{.ModelName}}Routes 获取{{.ModelName}}路由信息（用于集中注册） 
func Get{{.ModelName}}Routes() []RouteInfo {
    return []RouteInfo{
        {
            Method:      "POST",
            Path:        "Create{{.ModelName}}",
            Handler:     {{.ModelName | ToLower}}Api.Create{{.ModelName}},
            Description: "创建{{.ModelName}}",
        },
        {
            Method:      "GET", 
            Path:        "Get{{.ModelName}}List",
            Handler:     {{.ModelName | ToLower}}Api.Get{{.ModelName}}List,
            Description: "获取{{.ModelName}}列表",
        },
    }
}
`

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t, err := template.New("routes").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 创建目录
	//dir := outputPath[:strings.LastIndex(outputPath, "/")]
	dir := filepath.Dir(outputPath)
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
