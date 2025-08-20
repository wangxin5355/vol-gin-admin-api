package v1

import (
	"github.com/wangxin5355/vol-gin-admin-api/api/v1/example"
	"github.com/wangxin5355/vol-gin-admin-api/api/v1/system"
)

type ExampleApiGroup struct {
	example.ExampleTestApi
}

type SystemApiGroup struct {
	system.SystemApi
	system.AccountApi
}

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  SystemApiGroup
	ExampleApiGroup ExampleApiGroup
}
