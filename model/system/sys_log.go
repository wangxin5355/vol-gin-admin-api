package system
        
        
import (    "time"
)

// 系统日志 (Sys_Log)
type SysLog struct {
    //第一项是固定的 写描述信息
    _ struct{} `entity:"TableCnName=系统日志;TableName=Sys_Log;DetailTable=;DetailTableCnName=;DBServer=gin;Key=Id"`
    EndDate *time.Time `json:"EndDate" gorm:"column:EndDate;comment:结束时间"`
    Role_Id int `json:"Role_Id" gorm:"column:Role_Id;comment:角色ID"`
    User_Id int `json:"User_Id" gorm:"column:User_Id;comment:用户ID"`
    BrowserType *string `json:"BrowserType" gorm:"column:BrowserType;comment:浏览器类型"`
    ServiceIP *string `json:"ServiceIP" gorm:"column:ServiceIP;comment:服务器IP"`
    UserIP *string `json:"UserIP" gorm:"column:UserIP;comment:用户IP"`
    ExceptionInfo *string `json:"ExceptionInfo" gorm:"column:ExceptionInfo;comment:异常信息"`
    ResponseParameter *string `json:"ResponseParameter" gorm:"column:ResponseParameter;comment:响应参数"`
    RequestParameter *string `json:"RequestParameter" gorm:"column:RequestParameter;comment:请求参数"`
    ElapsedTime int `json:"ElapsedTime" gorm:"column:ElapsedTime;comment:时长"`
    Success int `json:"Success" gorm:"column:Success;comment:响应状态"`
    LogType *string `json:"LogType" gorm:"column:LogType;comment:日志类型"`
    Url *string `json:"Url" gorm:"column:Url;comment:请求地址"`
    UserName *string `json:"UserName" gorm:"column:UserName;comment:用户名称"`
    BeginDate *time.Time `json:"BeginDate" gorm:"column:BeginDate;comment:开始时间"`
    Id int `json:"Id" gorm:"column:Id;primaryKey"`
}

func (SysLog) TableName() string {
	return "Sys_Log"
}