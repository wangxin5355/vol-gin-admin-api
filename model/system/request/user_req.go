package request

// Register User register structure
type Register struct {
	Username  string `json:"userName" example:"用户名"`
	Password  string `json:"passWord" example:"密码"`
	NickName  string `json:"nickName" example:"昵称"`
	HeaderImg string `json:"headerImg" example:"头像链接"`
	Enable    int    `json:"enable" swaggertype:"string" example:"int 是否启用"`
	Phone     string `json:"phone" example:"电话号码"`
	Email     string `json:"email" example:"电子邮箱"`
}

// Login User login structure
type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}
