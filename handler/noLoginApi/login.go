package noLoginApi

import (
	"aDi/config"
	"aDi/dao"
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v2"
)

// Login 登录
func (n *NoLoginHandlerImp) Login(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.LoginReq{}
	defer util.TimeCost("Login", req, &rsp.Code)()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if req.Code == "" {
		log.Errorf("code is empty,code:%s", req.Code)
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	loginRsp, err := weapp.Login(config.GetAppId(), config.GetAppSecret(), req.Code)
	if err != nil {
		log.Errorf("login fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIS2S)
		return
	}
	if err = loginRsp.GetResponseError(); err != nil {
		log.Errorf("login fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrInvalidToken)
		return
	}
	log.Debugf("wx login,rsp:%+v", loginRsp)
	// 获取到对应的用户信息 -- 如果用户未注册则走注册流程
	userInfo, err := dao.GetUserInfoByOpenID(loginRsp.OpenID)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get user info fail,err:%s", err.Error())
		return
	}

	// 如果用户信息不存在 -- 则插入未注册表
	if userInfo == nil || userInfo.Uid <= 0 {
		unRegisterId, err := dao.UpsertUnRegisterInfo(loginRsp.OpenID, loginRsp.UnionID, loginRsp.SessionKey)
		if err != nil {
			log.Errorf("add unregister fail,err:%s", err.Error())
			rsp.WriteMsg(model.ErrIDbFail)
			return
		}
		// 返回未注册id
		rsp.Data = &model.LoginRsp{UnRegisterId: unRegisterId, ApiKey: service.GetSignApiKey(loginRsp.OpenID)}
		return
	}
	// 更新用户的session_key
	_ = dao.UpdateUserSessionKey(loginRsp.SessionKey, userInfo.Uid)
	// 如果用户信息存在 -- 则返回用户信息,并且返回对应的key code以及token
	tokenStr, expireAt, err := service.GenerateToken(userInfo.Uid)
	if err != nil {
		log.Errorf("generate token fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIInner)
		return
	}
	rsp.Data = &model.LoginRsp{
		ApiKey:    service.GetSignApiKey(userInfo.OpenID),
		User:      userInfo,
		TokenInfo: &model.TokenInfo{Token: tokenStr, ExpireAt: expireAt},
	}
	return
}

// Register 注册接口
func (n *NoLoginHandlerImp) Register(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.RegisterReq{}
	defer util.TimeCost("Register", req, &rsp.Code)()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if req.UnRegisterId <= 0 || req.WxDynamicCode == "" {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 先根据未注册id获取用户的未注册信息 -- 目前主要是用户的open id
	unRegisterInfo, err := dao.GetUnregisterInfoById(req.UnRegisterId)
	if err != nil {
		log.Errorf("get unregister info fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIDbFail)
		return
	}
	if unRegisterInfo == nil || unRegisterInfo.Id <= 0 {
		log.Error("用户未进行静默登录，需要先进行静默登录才能注册")
		rsp.WriteMsg(model.ErrINoLogin)
		return
	}
	// 调用微信接口获取手机号
	phone, err := service.GetPhoneByWxCode(req.WxDynamicCode)
	if err != nil {
		log.Errorf("get phone info fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIS2S)
		return
	}
	if phone == "" {
		rsp.WriteMsg(model.CodeMsg{Code: model.ErrIS2S.Code, Msg: "未获取到正确手机号,请重试"})
		return
	}
	// 进行用户注册 -- 向用户表中插入数据
	userInfo := &model.DBUserInfo{
		Nick:       "未命名",
		Age:        20,
		Sex:        "男",
		OpenID:     unRegisterInfo.OpenID,
		Phone:      phone,
		SessionKey: unRegisterInfo.SessionKey,
		UnionID:    unRegisterInfo.UnionID,
	}
	err = dao.AddUserInfo(userInfo)
	if err != nil {
		log.Errorf("add user info fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIDbFail)
		return
	}
	// 如果用户信息存在 -- 则返回用户信息,并且返回对应的key code以及token
	tokenStr, expireAt, err := service.GenerateToken(userInfo.Uid)
	if err != nil {
		log.Errorf("generate token fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIInner)
		return
	}
	rsp.Data = &model.RegisterRsp{
		User:      userInfo,
		TokenInfo: &model.TokenInfo{Token: tokenStr, ExpireAt: expireAt},
	}
	return rsp
}
