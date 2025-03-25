package noLoginApi

import (
	"aDi/config"
	"aDi/log"
	"github.com/medivhzhan/weapp/v2"
	"testing"
)

func TestLogin(t *testing.T) {
	loginRsp, err := weapp.Login(config.GetAppId(), config.GetAppSecret(), "0d3RtXZv33JmC43txm2w3l804p0RtXZf")
	if err != nil {
		log.Errorf("login fail,err:%s", err.Error())
		return
	}
	log.Infof("login rsp:%+v", loginRsp)
}
