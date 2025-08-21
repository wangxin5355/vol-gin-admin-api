package system

import (
	"time"
)

type Login interface {
	GetRoleIds() string
	GetDeptIds() string
	GetUserId() uint32
	GetUserInfo() any
	GetUsername() string
}

var _ Login = new(SysUser)

type SysUser struct {
	User_Id            uint32    `json:"user_Id" gorm:"primarykey" `
	Role_Ids           string    `json:"role_Ids" gorm:"comment:用户角色"`
	RoleName           string    `json:"RoleName" gorm:"column:RoleName;default:无;comment:角色名称，废弃"`
	PhoneNo            string    `json:"phoneNo" gorm:"comment:电话"`
	Remark             string    `json:"remark" gorm:"comment:备注"`
	Tel                string    `json:"tel" gorm:"comment:电话"`
	UserName           string    `json:"userName" gorm:"comment:用户名"`
	UserPwd            string    `json:"-"  gorm:"comment:用户登录密码"`
	UserTrueName       string    `json:"userTrueName" gorm:"comment:用户真是名称"`
	DeptName           string    `json:"deptName" gorm:"comment:部门名称"`
	Dept_Id            uint32    `json:"dept_Id" gorm:"comment:部门id"` // 用户登录密码
	Email              string    `json:"email" gorm:"comment:邮件"`
	Enable             int       `json:"enable" gorm:"comment:是否可用"`
	Gender             int       `json:"gender" gorm:"comment:性别"`
	HeadImageUrl       string    `json:"headImageUrl" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`
	IsRegregisterPhone int       `json:"isRegregisterPhone" gorm:"comment:是否注册手机"`
	LastLoginDate      time.Time `json:"lastLoginDate"  gorm:"comment:最后登录时间"`
	LastModifyPwdDate  time.Time `json:"lastModifyPwdDate"  gorm:"comment:最后修改密码时间"`
	Address            string    `json:"address" gorm:"comment:地址"`
	AppType            int       `json:"appType" gorm:"comment:app类型"`
	AuditDate          time.Time `json:"auditDate"  gorm:"comment:审核时间"`
	AuditStatus        int       `json:"auditStatus"  gorm:"comment:审核状态"`
	Auditor            string    `json:"auditor"  gorm:"comment:审核人名称"`
	OrderNo            int       `json:"orderNo"  gorm:"comment:排序号？"`
	Token              string    `json:"-"  gorm:"comment:用户token"`
	CreateID           int       `json:"createID"  gorm:"comment:创建人"`
	CreateDate         time.Time `json:"createDate"  gorm:"comment:创建时间"`
	Creator            string    `json:"creator"  gorm:"comment:创建人名称？"`
	Mobile             string    `json:"mobile"  gorm:"comment:手机号"`
	Modifier           string    `json:"modifier"  gorm:"comment:修改人名称？"`
	ModifyDate         time.Time `json:"modifyDate"  gorm:"comment:修改时间"`
	ModifyID           int       `json:"modifyID"  gorm:"comment:修改人id"`
	DeptIds            string    `json:"deptIds"  gorm:"comment:组织架构"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

func (s *SysUser) GetRoleIds() string {
	return s.Role_Ids
}

func (s *SysUser) GetDeptIds() string {
	return s.DeptIds
}

func (s *SysUser) GetUsername() string {
	return s.UserName
}

func (s *SysUser) GetUserId() uint32 {
	return s.User_Id
}

func (s *SysUser) GetUserInfo() any {
	return *s
}
