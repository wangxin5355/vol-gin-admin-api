package service

import (
	"github.com/wangxin5355/vol-gin-admin-api/service/example"
	"github.com/wangxin5355/vol-gin-admin-api/service/system"
	"github.com/wangxin5355/vol-gin-admin-api/service/test"
)

type TestServiceGroup struct {
	test.TestService
}
type ExampleServiceGroup struct {
	example.ExampleTestService
}
type SystemServiceGroup struct {
	system.OperationRecordService
	system.UserService
	system.SystemConfigService
	system.JwtService
	system.CasbinService
}

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  SystemServiceGroup
	ExampleServiceGroup ExampleServiceGroup
	TestServiceGroup    TestServiceGroup
}
