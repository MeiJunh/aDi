package model

type CosRegion string

const (
	CosSH = CosRegion("ap-shanghai")
)

// CosConfig cos配置信息
type CosConfig struct {
	AccessKey      string `json:"accessKey"`
	AccessSecret   string `json:"accessSecret"`
	BucketName     string `json:"bucketName"`
	Endpoint       string `json:"endpoint"`
	ResourceFormat string `json:"resourceFormat"`
	RoleArn        string `json:"roleArn"`
	CDNHost        string `json:"cdnHost"`
	HostName       string `json:"hostName"`
}

// OssSTSInfo sts回参信息
type OssSTSInfo struct {
	SecurityToken   string   `json:"securityToken"`   // token
	AccessKeyID     string   `json:"accessKeyID"`     // key
	AccessKeySecret string   `json:"accessKeySecret"` // secret密钥
	ExpireTime      int64    `json:"expireTime"`      // 过期时间 -- 单位秒,比如一个小时
	AvailDirList    []string `json:"availDirList"`    // 有权限的目录列表
	CDNHost         string   `json:"cdnHost"`         // cdn域名
}
