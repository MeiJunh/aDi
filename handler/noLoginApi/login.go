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
	defer util.TimeCost("login", req, &rsp.Code)()
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
		rsp.Data = &model.LoginRsp{UnRegisterId: unRegisterId}
		return
	}
	// 如果用户信息存在 -- 则返回用户信息,并且返回对应的key code以及token
	tokenStr, expireAt, err := service.GenerateToken(userInfo.Uid)
	if err != nil {
		log.Errorf("generate token fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIInner)
		return
	}
	rsp.Data = &model.LoginRsp{
		KeyCode:   service.GetSignAppKey(userInfo.OpenID),
		User:      userInfo,
		TokenInfo: &model.TokenInfo{Token: tokenStr, ExpireAt: expireAt},
	}
	return
}
