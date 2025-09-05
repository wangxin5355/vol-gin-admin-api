package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wangxin5355/vol-gin-admin-api/docs"
	"github.com/wangxin5355/vol-gin-admin-api/global"
	"github.com/wangxin5355/vol-gin-admin-api/middleware"
	"github.com/wangxin5355/vol-gin-admin-api/router"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	//TODO 初始化数据库
	//TODO 初始化redis
	Router := initRouters()
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))
	fmt.Printf("默认自动化文档地址:http://localhost%s/swagger/index.html", address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
}

func initRouters() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}
	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example
	testRouter := router.RouterGroupApp.Test

	//跨域设置
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	global.GVA_LOG.Info("use middleware cors")
	//Swagger
	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//统一加前缀
	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)  //不需要鉴权
	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix) //带鉴权
	PrivateGroup.Use(middleware.JWTAuth())                              //需要鉴权的路由
	// 健康监测
	PublicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	systemRouter.InitJwtRouter(PrivateGroup)
	{
		systemRouter.InitAccRouter(PublicGroup) // 注册Login,register 不需要鉴权的接口
	}
	{
		//传入PrivateGroup 则使用了鉴权
		//exampleRouter.InitTestRouter(PrivateGroup )
		//PublicGroup则不需要鉴权
		exampleRouter.InitTestRouter(PublicGroup)
		testRouter.InitTestRouter(PublicGroup)
		systemRouter.InitPermissionRouter(PublicGroup)
		systemRouter.InitMenuRouter(PrivateGroup) //需要鉴权
	}

	global.GVA_ROUTERS = Router.Routes()
	global.GVA_LOG.Info("router register success")
	return Router
}
