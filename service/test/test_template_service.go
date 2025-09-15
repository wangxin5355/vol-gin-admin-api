package test

import (
	"github.com/wangxin5355/vol-gin-admin-api/core/initialize"

	"github.com/wangxin5355/vol-gin-admin-api/core/base"
	"github.com/wangxin5355/vol-gin-admin-api/model/test"
	"github.com/wangxin5355/vol-gin-admin-api/model/test/partial"
)

func NewTestTemplateService() *TestTemplateService {
	return &TestTemplateService{
		BaseService: base.NewBaseService[partial.TestTemplateEntity, test.TestTemplate](string(initialize.DbGin)),
	}
}

// TestTemplateService 继承 BaseService[TestTemplate]
type TestTemplateService struct {
	*base.BaseService[partial.TestTemplateEntity, test.TestTemplate]
}
