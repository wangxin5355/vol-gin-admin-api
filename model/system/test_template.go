package system

type TestTemplate struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func (TestTemplate) TableName() string {
	return "test_template"
}
