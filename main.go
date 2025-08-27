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

	// 最后初始化服务，确保数据库已初始化
	service.InitServiceGroup()
}
