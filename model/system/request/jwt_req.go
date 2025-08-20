package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint64
	Username    string
	NickName    string
	AuthorityId uint
	TokenId     string
}
