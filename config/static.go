package config

import (
	"aDi/log"
	"aDi/model"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
)

// StaticConf 静态配置
type StaticConf struct {
	StaticDBDsn string                               `json:"staticDbDsn"` // 静态db dsn配置
	CosConfMap  map[model.CosRegion]*model.CosConfig `json:"cosConfMap"`
	WxPayConf   WxPayConf                            `json:"wxPayConf"` // 微信支付配置
}

// InitStaticConf 初始化静态配置信息
func InitStaticConf() {
	// 初始化环境信息
	configPath := "../conf/config.json"

	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("init config fail,%s\n", err.Error()))
		}
	}()

	var f *os.File
	f, err = os.Open(configPath)
	if err != nil {
		log.Errorf("InitConf open file fail,conf path:%s,err:%s", configPath, err.Error())
		return
	}
	var confByte []byte
	confByte, err = io.ReadAll(f)
	if err != nil {
		log.Errorf("InitConf ReadAll file fail,conf path:%s,err:%s", configPath, err.Error())
		return
	}

	err = jsoniter.Unmarshal(confByte, SConf)
	if err != nil {
		log.Errorf("InitConf Unmarshal file fail,conf path:%s,err:%s", configPath, err.Error())
		return
	}
	log.Infof("InitConf success,conf path:%s,config:%s", configPath, string(confByte))
}
