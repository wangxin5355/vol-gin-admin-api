package system

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"gorm.io/gorm"
	"strconv"
)

type UserService struct{}

//@author: wangxin
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.UserName).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	//生成六位随机包含数字和字母的密码
	newPwd := utils.GenerateRandomNumber(6)
	key := global.GVA_CONFIG.Secret.User
	pwd, err := utils.EncryptAES(newPwd, key)
	if err != nil {
		return userInter, err
	}
	u.UserPwd = pwd
	u.RoleName = "无"
	err = global.GVA_DB.Create(&u).Error
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

func (userService *UserService) GetUserInfoByCache(userId int32) (userInter *system.SysUser, err error) {

	//先从Redis拿，没有再从db获取
	userCacheKey := GetUserCacheKey(int(userId))
	user, err := utils.GetRedisStruct[system.SysUser](context.Background(), global.GVA_REDIS, userCacheKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if user.User_Id > 0 {
		return &user, nil
	}
	err = global.GVA_DB.Where("User_Id = ?", userId).First(&user).Error
	if err == nil {
		if user.User_Id > 0 {
			utils.SetRedisStruct(context.Background(), global.GVA_REDIS, userCacheKey, user, 0)
		}
		return &user, err
	} else {
		return nil, err
	}
}

func GetUserCacheKey(userId int) string {
	return model.UID.String() + strconv.Itoa(userId)
}
