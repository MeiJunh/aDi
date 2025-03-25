package loginApi

import (
	"aDi/dao"
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

/*
数字分身相关接口api
*/

// GetDigitalShadowInfo 获取自己数字分身信息
func (l *LoginHandlerImp) GetDigitalShadowInfo(c *gin.Context) (rsp model.BaseRsp) {
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetDigitalShadowInfo", uid, &rsp.Code)()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取个人数字分身信息
	digInfo, errCode, errMsg := service.GetDigitalShadowInfoByUid(uid)
	if errCode != model.ErrCodeSuccess {
		rsp.WriteMsg(model.CodeMsg{Code: errCode, Msg: errMsg})
		log.Errorf("get digital shadow info errCode:%d errMsg:%s", errCode, errMsg)
		return
	}
	rsp.Data = digInfo
	return
}

// AddDigitalShadow 添加数字分身
func (l *LoginHandlerImp) AddDigitalShadow(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.DigitalShadowInfo{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("AddDigitalShadow", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	// 获取用户信息 -- 填充用户基本信息
	// 新增更新数字人信息
	err = dao.UpsertDigitalInfo(uid, req)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("upsert digital info fail,err:%s", err.Error())
		return
	}
	return
}

// ModDigitalAllKinds 设置百变分身
func (l *LoginHandlerImp) ModDigitalAllKinds(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.DigitalAllKinds{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ModDigitalAllKinds", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	// 更新百变分身信息
	err = dao.UpdateDigitalField(uid, dao.DFDigitalAllKinds, util.MarshalToStringWithOutErr(req))
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("upsert digital info fail,err:%s", err.Error())
		return
	}
	return
}

// ModDigitalChargeInfo 设置百变分身收费
func (l *LoginHandlerImp) ModDigitalChargeInfo(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.ChargeInfoConf{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ModDigitalChargeInfo", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	// 更新百变分身收费
	err = dao.UpdateDigitalField(uid, dao.DFChargeConf, util.MarshalToStringWithOutErr(req))
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("upsert digital info fail,err:%s", err.Error())
		return
	}
	return
}
