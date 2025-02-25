package config

import "sync/atomic"

// GetSConfDsn 返回静态配置dsn
func GetSConfDsn() string {
	if SConf != nil {
		return SConf.StaticDBDsn
	}
	return ""
}

// DC 获取动态配置
func DC() *DynamicConf {
	p := atomic.LoadPointer(&mDcPt)
	return (*DynamicConf)(p)
}

func GetAppId() string {
	if DC() == nil {
		return ""
	}
	return DC().AppId
}

func GetAppSecret() string {
	if DC() == nil {
		return ""
	}
	return DC().AppSecret
}
