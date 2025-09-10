package partial

import "github.com/wangxin5355/vol-gin-admin-api/model/test"

// 测试模版 (test_template)
type TestTemplateEntity struct {
	test.TestTemplate
	//在这写自定义字段 例如 Test string `json:"test" gorm:"-"` //gorm:"-"表示忽略该字段
	Test string `json:"test" gorm:"-"` //gorm:"-"表示忽略该字段
}
