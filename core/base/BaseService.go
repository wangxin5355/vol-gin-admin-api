package base

import (
	"github.com/wangxin5355/vol-gin-admin-api/core/base/provider"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/request"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"gorm.io/gorm"
)

// BaseService 泛型基类
type BaseService[T any] struct {
	DB *gorm.DB
}

// 构造函数
func NewBaseService[T any](dbName string) *BaseService[T] {
	db := global.GetGlobalDBByDBName(dbName)
	if db == nil {
		panic("数据库连接未初始化或名称错误: " + dbName)
	}
	return &BaseService[T]{
		DB: db,
	}
}

// GetPageData 分页查询
func (s *BaseService[T]) GetPageData(options request.PageDataOptions) *response.PageGridData[T] {
	return provider.GetPageData[T](s.DB, options)
}

// Add 添加
func (s *BaseService[T]) Add(saveModel request.SaveModel) *response.WebResponseContent {
	return provider.Add[T](s.DB, saveModel)
}

// Update 更新
func (s *BaseService[T]) Update(saveModel request.SaveModel) *response.WebResponseContent {
	return provider.Update[T](s.DB, saveModel)
}

// Del 删除
func (s *BaseService[T]) Del(keys []any) *response.WebResponseContent {
	return provider.Del[T](s.DB, keys)
}
