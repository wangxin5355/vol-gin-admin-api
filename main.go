package main

import (
	"github.com/wangxin5355/vol-gin-admin-api/core"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"go.uber.org/zap"
)

// @title                       Gin-cli Swagger API接口文档
// @version                     v1.0.0
// @description                 使用gin快速开api脚手架
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	//初始化Viper 动态配置管理
	global.GVA_VP = core.Viper() // 初始化Viper
	//初始化Zap  日志
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	//初始化Redis
	initialize.Redis()
	//初始化gorm
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	//其他初始化检查
	initialize.OtherInit()

	if global.GVA_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	core.RunServer()
}
