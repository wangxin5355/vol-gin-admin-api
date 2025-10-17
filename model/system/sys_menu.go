package system

import (
	"github.com/wangxin5355/vol-gin-admin-api/model/dto"
	"time"
)

type SysMenu struct {
	Menu_Id     int          `json:"menu_Id" gorm:"primarykey" `
	MenuName    string       `json:"menuName" gorm:"comment: 菜单名称"`
	Auth        string       `json:"auth" gorm:"comment:可用菜单"`
	Icon        string       `json:"icon" gorm:"comment:图标Icon"`
	Description string       `json:"description" gorm:"comment:描述"`
	Enable      int          `json:"enable" gorm:"comment:是否可用"`
	OrderNo     int          `json:"orderNo"  gorm:"comment:排序号"`
	ITableName  string       `json:"tableName" gorm:"column:TableName;comment:表/视图名称"`
	ParentId    int          `json:"parentId" gorm:"comment:父菜单id"`
	Url         string       `json:"url" gorm:"comment:路由地址"`
	CreateDate  time.Time    `json:"createDate"  gorm:"comment:创建时间"`
	Creator     string       `json:"creator"  gorm:"comment:创建人"`
	ModifyDate  time.Time    `json:"modifyDate" gorm:"comment:修改时间"`
	Modifier    string       `json:"modifier" gorm:"comment:修改人"`
	MenuType    int          `json:"menuType"  gorm:"comment:菜单类型"`
	Actions     []dto.Action `json:"actions" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
