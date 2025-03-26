package service

import (
	"aDi/config"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"errors"
	"fmt"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"sync"
	"time"
)

var (
	cosCliMap sync.Map
)

// getCosClient 获取cos client
func getCosClient(regionId model.CosRegion) (*sts.Client, error) {
	cli, ok := cosCliMap.Load(regionId)
	if ok {
		return cli.(*sts.Client), nil
	}

	cInfo := config.GetCosConfMap(regionId)
	if cInfo == nil {
		return nil, errors.New("regionId not found")
	}

	cli = sts.NewClient(cInfo.AccessKey, cInfo.AccessSecret, nil)
	if cli == nil {
		log.Errorf("get oss sts client fail, region:%s, accessKey:%s, accessSecret:%s", regionId, cInfo.AccessKey, cInfo.AccessSecret)
		return nil, errors.New("get oss sts client fail")
	}
	log.Infof("get oss sts client success, region:%s, accessKey:%s, accessSecret:%s", regionId, cInfo.AccessKey, cInfo.AccessSecret)
	cosCliMap.Store(regionId, cli)
	return cli.(*sts.Client), nil
}

// GetTencentSTSByRegion 获取腾讯cos的sts信息
func GetTencentSTSByRegion(region model.CosRegion, uid, durationSeconds int64) (result *sts.CredentialResult, cdnHost string, dirList []string, err error) {
	c, err := getCosClient(region)
	if err != nil {
		log.Errorf("get cos sts client fail, region:%s, err:%s", region, err.Error())
		return result, cdnHost, nil, err
	}
	if c == nil {
		log.Errorf("get cos sts client fail, region:%s", region)
		return result, cdnHost, nil, errors.New("get cos sts client fail")
	}
	cInfo := config.GetCosConfMap(region)
	cdnHost = cInfo.CDNHost
	dirList = []string{fmt.Sprintf("yq/%d/%s", uid, time.Now().Format(util.DayNumFormat))}
	policy := GetTencentCosPoliceByRegion(cInfo.ResourceFormat, dirList)
	opt := &sts.CredentialOptions{
		DurationSeconds: durationSeconds,
		Region:          string(region),
		Policy:          policy,
		RoleSessionName: "yq",
	}
	result, err = c.GetCredential(opt)
	if err != nil {
		log.Errorf("get cos sts fail, region:%s, err:%s", region, err.Error())
		return result, cdnHost, nil, err
	}
	log.Infof("get cos sts success, region:%s, requestId:%s, credentials:%+v", region, result.RequestId, result.Credentials)
	return result, cdnHost, dirList, nil
}

func GetTencentCosPoliceByRegion(resourceFormat string, dirList []string) *sts.CredentialPolicy {
	resourceList := make([]string, 0, len(dirList))

	for _, v := range dirList {
		resourceList = append(resourceList, fmt.Sprintf(resourceFormat, v))
	}

	return &sts.CredentialPolicy{
		Statement: []sts.CredentialPolicyStatement{
			{
				Action: []string{
					"*",
				},
				Effect:   "allow",
				Resource: resourceList,
			},
		},
	}
}
