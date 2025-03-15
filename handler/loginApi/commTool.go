package loginApi

import (
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"github.com/gin-gonic/gin"
)

// GetCosSts 获取腾讯云上传sts信息
func (l *LoginHandlerImp) GetCosSts(c *gin.Context) (rsp model.BaseRsp) {
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetCosSts", uid, &rsp.Code)()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	expireTime := int64(util.HourSecond)
	cosInfo, cdnHost, dirList, err := service.GetTencentSTSByRegion(model.CosGZ, uid, expireTime)
	if err != nil {
		rsp.WriteMsg(model.ErrIS2S)
		log.Errorf("get cos sts err:%s", err.Error())
		return
	}
	rsp.Data = &model.OssSTSInfo{
		SecurityToken:   cosInfo.Credentials.SessionToken,
		AccessKeyID:     cosInfo.Credentials.TmpSecretID,
		AccessKeySecret: cosInfo.Credentials.TmpSecretKey,
		ExpireTime:      expireTime,
		AvailDirList:    dirList,
		CDNHost:         cdnHost,
	}
	return
}
