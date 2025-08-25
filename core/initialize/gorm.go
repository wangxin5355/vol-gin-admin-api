package initialize

import (
	"os"

	"github.com/wangxin5355/vol-gin-admin-api/config"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	case "mssql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mssql.Dbname
		return GormMssql()
	default:
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(
		system.SysUser{},
		system.SysOperationRecord{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}

// InitAllDB 初始化所有数据库连接
func InitAllDB() map[string]*gorm.DB {
	dbList := make(map[string]*gorm.DB)
	for _, dbConfig := range global.GVA_CONFIG.DBList {
		if dbConfig.Disable {
			continue
		}
		var db *gorm.DB
		switch dbConfig.Type {
		case "mysql":
			db = GormMysqlByConfig(config.Mysql{GeneralDB: dbConfig.GeneralDB})
		case "mssql":
			db = GormMssqlByConfig(config.Mssql{GeneralDB: dbConfig.GeneralDB})
		default:
			global.GVA_LOG.Warn("不支持的数据库类型", zap.String("type", dbConfig.Type))
			continue
		}
		dbList[dbConfig.AliasName] = db
	}
	global.GVA_DBList = dbList
	return dbList
}

// 定义数据库类型枚举
type DbTypeEnum string

const (
	DbGin DbTypeEnum = "gin"
)
