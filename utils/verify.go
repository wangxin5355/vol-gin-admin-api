package utils

var (
	LoginVerify               = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify            = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}}
	UpdateUserRoleVerify      = Rules{"UserId": {Gt("0")}, "RoleIds": {NotEmpty()}}
	CheckRolePermissionVerify = Rules{"RoleId": {Gt("0")}, "MenuId": {Gt("0")}, "Action": {NotEmpty()}}
	MenuMetaVerify            = Rules{"Title": {NotEmpty()}}

	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
)
