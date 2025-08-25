package initialize

import (
	"github.com/wangxin5355/vol-gin-admin-api/config"
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize/internal"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// GormMssql 初始化Mssql数据库
// Author [LouisZhang](191180776@qq.com)
func GormMssql() *gorm.DB {
	m := global.GVA_CONFIG.Mssql
	if m.Dbname == "" {
		return nil
	}
	mssqlConfig := sqlserver.Config{
		DSN:               m.Dsn(), // DSN data source name
		DefaultStringSize: 191,     // string 类型字段的默认长度
	}
	if db, err := gorm.Open(sqlserver.New(mssqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// GormMssqlByConfig 初始化Mysql数据库用过传入配置
func GormMssqlByConfig(m config.Mssql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mssqlConfig := sqlserver.Config{
		DSN:               m.Dsn(), // DSN data source name
		DefaultStringSize: 191,     // string 类型字段的默认长度
	}
	if db, err := gorm.Open(sqlserver.New(mssqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// GormMssqlWithConfig 支持根据 GeneralDB 初始化 MSSQL 数据库连接
func GormMssqlWithConfig(general config.GeneralDB) (*gorm.DB, error) {
	if general.Dbname == "" {
		return nil, nil
	}
	dsn := "sqlserver://" + general.Username + ":" + general.Password + "@" + general.Path + ":" + general.Port + "?database=" + general.Dbname + "&" + general.Config
	mssqlConfig := sqlserver.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}
	db, err := gorm.Open(sqlserver.New(mssqlConfig), internal.Gorm.Config(general.Prefix, general.Singular))
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(general.MaxIdleConns)
	sqlDB.SetMaxOpenConns(general.MaxOpenConns)
	return db, nil
}
