package router

import (
	"github.com/wangxin5355/vol-gin-admin-api/router/example"
	"github.com/wangxin5355/vol-gin-admin-api/router/system"
	"github.com/wangxin5355/vol-gin-admin-api/router/test"
)

type TestRouterGroup struct {
	test.TestRouter
}
type ExampleRouterGroup struct {
	example.ExampleTestRouter
}

type SystemRouterGroup struct {
	system.JwtRouter
	system.AccountRouter
	system.PermissionRouter
	system.SysTableInfoRouter
}

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  SystemRouterGroup
	Example ExampleRouterGroup
	Test    TestRouterGroup
}
