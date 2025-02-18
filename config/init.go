package config

import (
	"aDi/log"
	"unsafe"
)

var (
	SConf = &StaticConf{
		StaticDBDsn: "", // -- 默认为开发环境db配置
	}
	mDConf = &DynamicConf{} // 动态配置信息 -- 通过数据库动态读取
	mDcPt  unsafe.Pointer
)

// DynamicConf 动态配置列表
type DynamicConf struct {
}

// Init 配置初始化
func Init() {
	// 初始化静态配置
	InitStaticConf()
	// 初始化动态配置
	InitDynamicConf()
	return
}

// InitDynamicConf 初始化动态配置信息
func InitDynamicConf() {
	source, err := NewSQLConfSourceByURL(GetSConfDsn())
	if err != nil {
		log.Errorf("new sql source fail,err:%s", err.Error())
		return
	}
	watchList := []*MCWatchInfo{
		{
			Def: mDConf,
			PT:  &mDcPt,
			Key: "dynamic-key", // config对应的key为dynamic-key
		},
	}
	// service name设置为comm
	mc, err := NewConfig(source, watchList, AddServiceName("aDi"))
	if err != nil {
		log.Errorf("new config fail,err:%s", err.Error())
		return
	}

	mc.Watch()
}
