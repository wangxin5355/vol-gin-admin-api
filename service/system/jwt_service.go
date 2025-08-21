package system

import (
	"context"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/utils"
	"strconv"
)

type JwtService struct{}

var JwtServiceApp = new(JwtService)

//@author: wangxin
//@function: AddInBlacklist
//@description: 拉黑jwt
//@param: string tokenid
//@return: err error

func (jwtService *JwtService) AddInBlacklist(jwt string) (err error) {
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.SAdd(context.Background(), "TokenBlacklist", jwt, timer).Err()
	return err
}

//@author: wanxgin
//@function: IsBlacklist
//@description: 判断 tokenid是否存在黑名单redis中
//@param: tokenid string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	exists, err := global.GVA_REDIS.SIsMember(context.Background(), "TokenBlacklist", jwt).Result()
	if err != nil {
		panic(err)
	}
	return exists

}

//@author: wangxin
//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (jwtService *JwtService) GetRedisJWT(userId int) (redisJWT string, err error) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), strconv.Itoa(userId)).Result()
	return redisJWT, err
}

//@author: wangxin
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (jwtService *JwtService) SetRedisJWT(jwt string, userId int) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), strconv.Itoa(userId), jwt, timer).Err()
	return err
}
