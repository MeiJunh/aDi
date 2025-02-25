package handler

import (
	"aDi/handler/loginApi"
	"aDi/handler/noLoginApi"
	"github.com/gin-gonic/gin"
)

// MCheckHandler handler实例
var MCheckHandler *loginApi.LoginHandlerImp
var MNoCheckHandler *noLoginApi.NoLoginHandlerImp

// Init 路由注册,服务启动
func Init() {
	// 模版实例初始化
	MCheckHandler = &loginApi.LoginHandlerImp{}
	MNoCheckHandler = &noLoginApi.NoLoginHandlerImp{}
	// 初始化路由
	InitRoute()
}

func InitRoute() {
	router := gin.New()
	// 自动路由
	router.Any("/api/*action", CheckDynamicRouter(MCheckHandler))
	router.Any("/apiN/*action", NoCheckDynamicRouter(MNoCheckHandler))
	router.Run(":18888")
}
