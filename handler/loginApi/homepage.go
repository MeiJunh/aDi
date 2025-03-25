package loginApi

import (
	"aDi/dao"
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

/*
个人主页相关
*/

// UpsertHomepageConf 个人主页信息设置
func (l *LoginHandlerImp) UpsertHomepageConf(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.HpConf{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("UpsertHomepageConf", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 新增或者更新个人主页配置
	err = dao.UpsertHomepageConf(uid, req)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("upsert homepage conf fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	return
}

// GetMyHomepageConf 个人主页信息获取
func (l *LoginHandlerImp) GetMyHomepageConf(c *gin.Context) (rsp model.BaseRsp) {
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyHomepageConf", uid, &rsp.Code)()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}

	// 获取个人主页配置
	hpConf, err := dao.GetHomepageConfByUid(uid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get homepage conf fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	if hpConf == nil || hpConf.ID <= 0 || hpConf.HpConf == nil {
		// 数据初始化
		rsp.Data = &model.HpConf{
			CardList: make([]*model.HpCard, 0),
		}
		return
	}
	rsp.Data = hpConf.HpConf
	return
}
