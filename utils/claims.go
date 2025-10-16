package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/model/system"
	systemReq "github.com/wangxin5355/vol-gin-admin-api/model/system/request"
	"net"
	"strings"
)

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	//if token == "" {
	//	j := NewJWT()
	//	token, _ = c.Cookie("x-token")
	//	claims, err := j.ParseToken(token)
	//	if err != nil {
	//		global.GVA_LOG.Error("重新写入cookie token失败,未能成功解析token,请检查请求头是否存在x-token且claims是否为规定结构")
	//		return token
	//	}
	//	SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
	//}
	return token
}

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint32 {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.UserID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.BaseClaims.UserID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
//func GetUserAuthorityId(c *gin.Context) uint {
//	if claims, exists := c.Get("claims"); !exists {
//		if cl, err := GetClaims(c); err != nil {
//			return 0
//		} else {
//			return cl.AuthorityId
//		}
//	} else {
//		waitUse := claims.(*systemReq.CustomClaims)
//		return waitUse.AuthorityId
//	}
//}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Username
	}
}

func LoginToken(user system.Login) (token string, claims systemReq.CustomClaims, err error) {
	j := NewJWT()
	claims = j.CreateClaims(systemReq.BaseClaims{
		UserID:   user.GetUserId(),
		Username: user.GetUsername(),
		Role_Ids: user.GetRoleIds(),
		DeptIds:  user.GetDeptIds(),
	})
	token, err = j.CreateToken(claims)
	return
}

// 获取菜单类型
func GetMenuType(c *gin.Context) int {
	uapp := c.Request.Header.Get("uapp")
	if len(uapp) == 0 {
		return 0
	} else {
		return 1
	}
}

// 获取用户roeids
func GetUserRoles(c *gin.Context) ([]int, error) {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil, err
		} else {
			//从claims拿roleIDs
			roleIds := strings.Split(cl.Role_Ids, ",")
			var roleIds_int = StringSliceToIntSliceFilter(roleIds)
			return roleIds_int, nil
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		roleIds := strings.Split(waitUse.Role_Ids, ",")
		var roleIds_int = StringSliceToIntSliceFilter(roleIds)
		return roleIds_int, nil
	}
}

// 获取用户roeids
func GetUserRolesStr(c *gin.Context) ([]string, error) {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil, err
		} else {
			//从claims拿roleIDs
			roleIds := strings.Split(cl.Role_Ids, ",")
			return roleIds, nil
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		roleIds := strings.Split(waitUse.Role_Ids, ",")
		return roleIds, nil
	}
}
