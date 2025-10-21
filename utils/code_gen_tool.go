package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GenerateTest() {
	var modelNames = []string{"User", "Order"}                     //后面通过反射/db里面代码生产数据拿到模型名称
	var outputPath = "E:\\学习资料\\GoDev\\vol-gin-admin-api\\router1" //后面通过代码获取当前项目路径
	for _, modelName := range modelNames {
		fmt.Println("开始生成:" + modelName + "路径：" + filepath.Join(outputPath, modelName))
		GenerateRouteFile(RouteConfig{ModelName: modelName, BasePath: strings.ToLower(modelName)}, filepath.Join(outputPath, strings.ToLower(modelName), strings.ToLower(modelName)+"_route.go"))
	}
	GenerateMainRouterFile(MainRouterConfig{Models: modelNames}, filepath.Join(outputPath, "route.go"))
}
