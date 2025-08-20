package system

import (
	"context"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
)

type JwtService struct{}

var JwtServiceApp = new(JwtService)

//@author: wangxin
//@function: AddInBlacklist
//@description: 拉黑jwt
//@param: string tokenid
//@return: err error

func (jwtService *JwtService) AddInBlacklist(tokenid string, jwt string) (err error) {
	//TODO 往RedisSet里面加
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), "TokenBlacklist:"+tokenid, jwt, timer).Err()
	return err
}

//@author: wanxgin
//@function: IsBlacklist
//@description: 判断 tokenid是否存在黑名单redis中
//@param: tokenid string
//@return: bool

func (jwtService *JwtService) IsBlacklist(tokenid string) bool {
	exists, err := global.GVA_REDIS.Exists(context.Background(), "TokenBlacklist:"+tokenid).Result()
	if err != nil {
		panic(err)
	}
	if exists == 1 {
		return true
	} else {
		return false
	}

}

//@author: wangxin
//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), userName).Result()
	return redisJWT, err
}

//@author: wangxin
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}
