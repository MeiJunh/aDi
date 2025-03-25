package config

import (
	"aDi/config/dynamic"
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
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	AiConf    AiConf `json:"aiConf"` // ai配置
}

// WxPayConf 微信支付配置
type WxPayConf struct {
	AppId                              string `json:"appId"`                              // 小程序id
	WxDirectMchId                      string `json:"WxDirectMchId,omitempty"`            // 直连商户微信支付商户ID
	WxDirectMchCertificateSerialNumber string `json:"WxDirectMchCertificateSerialNumber"` // 直连商户证书序列号
	WxDirectAPIv3Key                   string `json:"WxDirectAPIv3Key"`                   // 直连商户微信支付APIv3密钥
	WxDirectApiClientKey               string `json:"WxDirectApiClientKey"`               // 直连商户微信支付证书密钥pem格式（apiclient_key.pem）
	WxPublicKeyStr                     string `json:"WxPublicKeyStr"`                     // 直连商户微信支付证书微信支付公钥
	WxPublicKeyId                      string `json:"WxPublicKeyId"`                      // 公钥ID
	WxDirectPayNotify                  string `json:"WxDirectPayNotify"`                  // 微信支付回调地址
}

// AiConf ai相关配置
type AiConf struct {
	ApiUrl      string `json:"apiUrl"`      // 请求地址url
	Secret      string `json:"secret"`      // 接口请求对应的secret
	TextAiModel string `json:"textAiModel"` // 文本ai模型名
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
	source, err := dynamic.NewSQLConfSourceByURL(GetSConfDsn())
	if err != nil {
		log.Errorf("new sql source fail,err:%s", err.Error())
		return
	}
	watchList := []*dynamic.MCWatchInfo{
		{
			Def: mDConf,
			PT:  &mDcPt,
			Key: "dynamic-key", // config对应的key为dynamic-key
		},
	}
	// service name设置为comm
	mc, err := dynamic.NewConfig(source, watchList, dynamic.AddServiceName("aDi"))
	if err != nil {
		log.Errorf("new config fail,err:%s", err.Error())
		return
	}

	mc.Watch()
}
