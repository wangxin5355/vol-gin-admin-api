package response

import (
	"github.com/wangxin5355/vol-gin-admin-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, businesscode int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, WebResponseContent{
		Status:  code == SUCCESS, //成功失败状态，如果false打印错误信息
		Code:    businesscode,    //业务编码
		Message: msg,
		Data:    data,
	})
}

// 响应内容
type WebResponseContent struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Ok(msg string, data any) *WebResponseContent {
	if msg == "" {
		msg = "操作成功"
	}
	if data == nil {
		data = nil // 或 data = map[string]interface{}{}
	}
	return &WebResponseContent{Status: true, Message: msg, Data: data}
}

func Error(msg string) *WebResponseContent {
	return &WebResponseContent{Status: false, Message: msg}
}

// WebResponse 返回WebResponseContent内容
func WebResponse(resp *WebResponseContent, c *gin.Context) {
	c.JSON(http.StatusOK, resp)
}

func OkWithContext(c *gin.Context) {
	Result(SUCCESS, model.GeneralSuccess, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, model.GeneralSuccess, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, model.GeneralSuccess, data, "成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, model.GeneralSuccess, data, message, c)
}
func OkWithDetailedAndBusinessCode(businesscode int, data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, businesscode, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, model.GeneralError, map[string]interface{}{}, "操作失败", c)
}
func FailWithBusinessCode(businesscode int, c *gin.Context) {
	Result(ERROR, businesscode, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, model.GeneralError, map[string]interface{}{}, message, c)
}

func FailWithMessageAndBusinessCode(message string, businesscode int, c *gin.Context) {
	Result(ERROR, businesscode, map[string]interface{}{}, message, c)
}

func NoAuth(message string, c *gin.Context) {
	//c.JSON(http.StatusUnauthorized, Response{
	//	7,
	//	nil,
	//	message,
	//})
	Result(ERROR, model.NoPermissions, nil, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, model.GeneralError, data, message, c)
}

func FailWithDetailedAndBusinessCode(data interface{}, businesscode int, message string, c *gin.Context) {
	Result(ERROR, businesscode, data, message, c)
}
