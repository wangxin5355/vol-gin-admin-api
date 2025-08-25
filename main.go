//go:generate swag init
package main

import (
	"github.com/wangxin5355/vol-gin-admin-api/core"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"go.uber.org/zap"
)

// @title                       vol-gin-admin-api Swagger API接口文档
// @version                     v1.0.0
// @description                 vol快速开发框架基于go语言gin框架实现
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	core.RunServer()
}

func init() {
	//初始化Viper 动态配置管理
	global.GVA_VP = core.Viper() // 初始化Viper
	//初始化Zap  日志
	global.GVA_LOG = core.Zap() // 初始化zap日志库

	zap.ReplaceGlobals(global.GVA_LOG)

	//初始化gorm
	// 初始化多库
	initialize.InitAllDB()
	// 初始化主库,获取gorm.dblist中的第一个作为主库
	global.GVA_DB = global.GetFirstDB()
	//global.GVA_DB = initialize.Gorm()

	//初始化Redis
	initialize.Redis()
	//其他初始化检查
	initialize.OtherInit()

	// 初始化服务分组，确保数据库已初始化
	importedServiceInit()
}

// 新增一个辅助函数用于调用服务初始化，避免 import 循环
func importedServiceInit() {
	// 这里直接调用 service 包的初始化方法
	// 需要确保 service 包已正确导入
	// 如果 service 包未导入，请在 import 中添加
	// "github.com/wangxin5355/vol-gin-admin-api/service"
	service.InitServiceGroup()
}
