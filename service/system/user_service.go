package system

import (
	"errors"
	"fmt"

	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"gorm.io/gorm"
)

//@author: wangxin
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

type UserService struct{}

var UserServiceApp = new(UserService)

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.UserName).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	key := global.GVA_CONFIG.Secret.User
	pwd, err := utils.EncryptAES(u.UserPwd, key)
	if err != nil {
		return userInter, err
	}
	// if pwd != "j79rYYvCz4vdhcboB1Ausg==" {
	// 	return u, errors.New("密码错误")
	// }
	u.UserPwd = pwd
	u.RoleName = "无"
	//err = global.GVA_DB.Create(&u).Error
	return u, err
}

//@author: wangxin
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.GVA_DB {
		return nil, fmt.Errorf("db not init")
	}
	var user system.SysUser
	err = global.GVA_DB.Where("UserName = ?", u.UserName).First(&user).Error
	if err == nil {
		key := global.GVA_CONFIG.Secret.User
		if ok, _ := utils.DecryptAES(user.UserPwd, key); ok == u.UserPwd {
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}
