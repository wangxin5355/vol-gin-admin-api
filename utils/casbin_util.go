package utils

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"go.uber.org/zap"
	"sync"
)

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

// GetCasbin 获取casbin实例
func GetCasbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(global.GVA_DB)
		if err != nil {
			zap.L().Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
			return
		}

		m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act,eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
		if err != nil {
			zap.L().Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		syncedCachedEnforcer, err = casbin.NewSyncedCachedEnforcer(m, a)
		if err != nil {
			zap.L().Error(err.Error(), zap.Error(err))
			return
		}
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}
