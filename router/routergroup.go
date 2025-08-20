package router

import (
	"github.com/wangxin5355/vol-gin-admin-api/router/example"
	"github.com/wangxin5355/vol-gin-admin-api/router/system"
)

type ExampleRouterGroup struct {
	example.ExampleTestRouter
}

type SystemRouterGroup struct {
	system.JwtRouter
	system.AccountRouter
}

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  SystemRouterGroup
	Example ExampleRouterGroup
}
