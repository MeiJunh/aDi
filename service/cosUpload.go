package service

import (
	"aDi/config"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"sync"
)

var (
	cosUploadCliMap = sync.Map{}
)

// getCosUploadClient 获取cos文件上传client
func getCosUploadClient(regionId model.CosRegion) (*cos.Client, error) {
	cli, ok := cosUploadCliMap.Load(regionId)
	if ok {
		return cli.(*cos.Client), nil
	}

	cInfo := config.GetCosConfMap(regionId)
	if cInfo == nil {
		return nil, errors.New("regionId not found")
	}

	cli = InitCosClient(fmt.Sprintf("https://%s", cInfo.Endpoint), cInfo.AccessKey, cInfo.AccessSecret)
	if cli == nil {
		log.Errorf("get oss upload client fail, region:%s, accessKey:%s, accessSecret:%s", regionId, cInfo.AccessKey, cInfo.AccessSecret)
		return nil, errors.New("get oss sts client fail")
	}
	log.Infof("get oss upload client success, region:%s, accessKey:%s, accessSecret:%s", regionId, cInfo.AccessKey, cInfo.AccessSecret)
	cosUploadCliMap.Store(regionId, cli)
	return cli.(*cos.Client), nil
}

// InitCosClient 初始化cos 腾讯云client
func InitCosClient(rawUrl, secretID, secretKey string) (cosClient *cos.Client) {
	u, _ := url.Parse(rawUrl)
	b := &cos.BaseURL{BucketURL: u}
	cosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: secretID, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return cosClient
}

// CosDoUpload cos进行文件上传
func CosDoUpload(data []byte, objectName string, tencentOptions *cos.ObjectPutHeaderOptions) (url string, err error) {
	defer util.TimeCost("CosDoUpload", objectName, &url)()
	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions:       nil,
		ObjectPutHeaderOptions: tencentOptions,
	}
	reader := bytes.NewReader(data)
	cosCli, err := getCosUploadClient(model.CosSH)
	if err != nil {
		log.Errorf("get cos upload client fail, region:%s, err:%s", model.CosSH, err.Error())
		return "", err
	}
	_, err = cosCli.Object.Put(context.Background(), objectName, reader, opt)
	if err != nil {
		return "", err
	}
	return config.GetCosConfMap(model.CosSH).CDNHost + objectName, err
}
