package system

import (
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	systemReq "github.com/wangxin5355/vol-gin-admin-api/model/system/request"
	systemRes "github.com/wangxin5355/vol-gin-admin-api/model/system/response"
	"github.com/wangxin5355/vol-gin-admin-api/service"
	"github.com/wangxin5355/vol-gin-admin-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type AccountApi struct{}

var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService

// Login
// @Tags     AccountApi
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      systemReq.LoginReq  true  "用户名, 密码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResp,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /acc/login [post]
func (b *AccountApi) Login(c *gin.Context) {
	var l systemReq.LoginReq
	err := c.ShouldBindJSON(&l)
	//key := c.ClientIP()

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	u := &system.SysUser{UserName: l.Username, UserPwd: l.Password}
	user, err := userService.Login(u)
	if err != nil {
		global.GVA_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		response.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	if user.Enable != 1 {
		global.GVA_LOG.Error("登陆失败! 用户被禁止登录!")
		response.FailWithMessage("用户被禁止登录", c)
		return
	}
	b.TokenNext(c, *user)
	return

	response.FailWithMessage("验证码错误", c)
}

// TokenNext 登录以后签发jwt
func (b *AccountApi) TokenNext(c *gin.Context, user system.SysUser) {
	token, claims, err := utils.LoginToken(&user)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	//if !global.GVA_CONFIG.System.UseMultipoint {
	//	utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	//	response.OkWithDetailed(systemRes.LoginResponse{
	//		User:      user,
	//		Token:     token,
	//		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	//	}, "登录成功", c)
	//	return
	//}
	//先从redis查询已经存在的token，如果没有就直接存入新的，
	if jwtStr, err := jwtService.GetRedisJWT(int(user.User_Id)); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, int(user.User_Id)); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		//utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResp{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil { //取老token报错，直接失败
		global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		//老token放入黑名单
		if err := jwtService.AddInBlacklist(jwtStr); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, int(user.User_Id)); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResp{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

// Register
// @Tags     AccountApi
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body     systemReq.Register   true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResp,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /acc/register [post]
func (b *AccountApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user := &system.SysUser{UserName: r.Username, UserPwd: r.Password, HeadImageUrl: r.HeaderImg, Enable: r.Enable, Mobile: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResp{User: userReturn}, "注册失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysUserResp{User: userReturn}, "注册成功", c)
}
