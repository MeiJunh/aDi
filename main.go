package main

import (
	"aDi/config"
	"aDi/dao"
	"aDi/handler"
	"aDi/log"
	"aDi/service"
)

func main() {
	// 日志初始化
	log.Init(false, "../log/aDi.log")
	// 配置初始化
	config.Init()
	// db初始化
	dao.Init()
	// service初始化
	service.Init()
	// 路由注册
	handler.Init()
}
