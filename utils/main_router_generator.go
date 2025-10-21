// utils/main_router_generator.go
package utils

import (
	"os"
	"strings"
	"text/template"
)

// MainRouterConfig 主路由配置
type MainRouterConfig struct {
	Models []string //后面通过反射拿models
}

// GenerateMainRouterFile 生成主路由文件
func GenerateMainRouterFile(config MainRouterConfig, outputPath string) error {
	tmpl := `package router

import (
    "github.com/gin-gonic/gin"
    {{range $index, $model := .Models}}
    "github.com/wangxin5355/vol-gin-admin-api/routes/{{$model | ToLower}}"{{end}}
)

// SetupRoutes 设置所有路由
func SetupRoutes(router *gin.Engine) {
    // API路由组
    apiGroup := router.Group(global.GVA_CONFIG.System.RouterPrefix) //统一前缀

    // 注册各模块路由
    {{range .Models}}
    {{. | ToLower}}.Register{{.}}Routes(apiGroup)
    {{end}}

}

// GetAllRoutes 获取所有路由信息（用于文档或监控）
func GetAllRoutes() []RouteInfo {
    var allRoutes []RouteInfo
    {{range .Models}}
    // {{.}}路由
    allRoutes = append(allRoutes, {{. | ToLower}}.Get{{.}}Routes()...)
    {{end}}
    return allRoutes
}
`

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t, err := template.New("main_router").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, config)
}
