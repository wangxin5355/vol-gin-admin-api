package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

// CustomClaims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	ID       uint32 //用户id
	Username string
	Role_Ids string //角色集合
	DeptIds  string //部门集合
}
