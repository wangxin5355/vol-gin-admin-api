package model

type CPrefix int

const (
	Role     CPrefix = iota // 0
	UID                     // 1
	HDImg                   // 2
	Token                   // 3
	CityList                // 4
)

func (s CPrefix) String() string {
	switch s {
	case Role:
		return "0"
	case UID:
		return "1"
	case HDImg:
		return "2"
	case Token:
		return "3"
	case CityList:
		return "4"
	default:
		return "UNKNOWN"
	}
}

const (
	GeneralError      = -1
	GeneralSuccess    = 0 //通用成功码
	ServerError       = 1
	LoginExpiration   = 302 //和netcore版本保持一致
	ParametersLack    = 303
	TokenExpiration   = 304
	PINError          = 305
	NoPermissions     = 306
	NoRolePermissions = 307
	LoginError        = 308
	AccountLocked     = 309
	LoginSuccess      = 310 //登录成功
	SaveSuccess       = 311
	AuditSuccess      = 312
	OperSuccess       = 313
	RegisterSuccess   = 314
	ModifyPwdSuccess  = 315
	EidtSuccess       = 316
	DelSuccess        = 317
	NoKey             = 318
	NoKeyDel          = 319
	KeyError          = 320
	Other             = 321 // 4
)
