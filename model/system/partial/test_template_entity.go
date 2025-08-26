package partial

import "github.com/wangxin5355/vol-gin-admin-api/model/system"

type TestTemplateEntity struct {
	system.TestTemplate
	Test string `json:"test" gorm:"-"`
}
