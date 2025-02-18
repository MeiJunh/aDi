package handler

import (
	"github.com/gin-gonic/gin"
)

// MCheckHandlerImp handler模版 -- 需要进行cookie校验
type MCheckHandlerImp struct {
}

// MNoCheckHandlerImp handler模版 -- 不需要进行校验
type MNoCheckHandlerImp struct {
}

// MCheckHandler handler实例
var MCheckHandler *MCheckHandlerImp
var MNoCheckHandler *MNoCheckHandlerImp

// Init 路由注册,服务启动
func Init() {
	// 模版实例初始化
	MCheckHandler = &MCheckHandlerImp{}
	MNoCheckHandler = &MNoCheckHandlerImp{}
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
