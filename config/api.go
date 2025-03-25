package config

import (
	"aDi/model"
	"sync/atomic"
)

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

// GetCosConfMap 获取cos conf
func GetCosConfMap(region model.CosRegion) (conf *model.CosConfig) {
	if SConf == nil || SConf.CosConfMap == nil {
		return conf
	}
	return SConf.CosConfMap[region]
}

func GetAiApiUrl() string {
	if DC() == nil {
		return ""
	}
	return DC().AiConf.ApiUrl
}

func GetAiSecret() string {
	if DC() == nil {
		return ""
	}
	return DC().AiConf.Secret
}

func GetAiTextAiModel() string {
	if DC() == nil {
		return ""
	}
	return DC().AiConf.TextAiModel
}

// GetWxPayConf 获取微信支付配置
func GetWxPayConf() WxPayConf {
	if SConf == nil {
		return WxPayConf{}
	}
	return SConf.WxPayConf
}
