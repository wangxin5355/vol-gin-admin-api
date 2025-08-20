package example

import (
	"github.com/gin-gonic/gin"
	"github.com/wangxin5355/vol-gin-admin-api/model/common/response"
)

type ExampleTestApi struct{}

// Test
// @Tags     ExampleTest
// @Summary  测试1
// @Produce   application/json
// @Param    id  query  int  true  "用户ID"
// @Success  200    {object}  response.Response{data=string,msg=string} "获取meta信息"
// @Router   /example/Test [get]
func (b *ExampleTestApi) Test(c *gin.Context) {
	key := c.ClientIP()
	response.OkWithMessage(key, c)
}

// Test1
// @Tags     ExampleTest
// @Summary  测试2
// @Produce   application/json
// @Param    id  query  int  true  "用户ID"
// @Success  200   {object}  response.Response{data=string,msg=string} "获取meta信息"
// @Router    /example/Test1 [get]
func (b *ExampleTestApi) Test1(c *gin.Context) {
	key := c.ClientIP()
	response.OkWithMessage(key, c)
}
