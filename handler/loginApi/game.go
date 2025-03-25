package loginApi

import (
	"aDi/dao"
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetMyGameList 获取我的游戏列表 -- 用户获取自己的游戏列表
func (l *LoginHandlerImp) GetMyGameList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyGameList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	// 获取用户的游戏配置列表
	gameList, hasMore, nextIndex, err := dao.GetMyGameList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get my game list fail,err:%s", err.Error())
		return
	}

	rsp.Data = &model.GetListRsp{
		List:      model.TransDbGameToGameInfo(gameList),
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// CreateGame 创建游戏
func (l *LoginHandlerImp) CreateGame(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GameInfo{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("CreateGame", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	// 参数校验
	if req.RETotalAmount <= 0 || req.RETotalNum <= 0 || req.Prologue == "" || len(req.AnswerList) <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取用户信息
	userInfo, err := dao.GetUserInfoByUid(uid)
	if err != nil {
		log.Errorf("get user info fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIDbFail)
		return
	}
	if userInfo == nil || userInfo.Uid <= 0 {
		rsp.WriteMsg(model.CodeMsg{Code: model.ECNoExist, Msg: "该用户未注册"})
		return
	}
	// 创建订单 -- 并且记录,返回订单信息
	orderInfo, err := service.WxMchJsapi(context.Background(), &model.UnifiedOrderReq{
		Uid:                uid,
		OpenId:             userInfo.OpenID,
		ProductDescription: "创建五句话挑战游戏",
		Amount:             req.RETotalAmount,
		ProdType:           model.PTRe,
		ExpandStr:          util.MarshalToStringWithOutErr(req),
	})
	if err != nil {
		log.Errorf("add pay order fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIInner)
		return
	}
	rsp.Data = orderInfo
	return
}

// GetGamePlayListByShareCode 根据分享码获取游戏对话信息列表
func (l *LoginHandlerImp) GetGamePlayListByShareCode(c *gin.Context) (rsp model.BaseRsp) {
	shareCode := c.Query("shareCode")
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetGamePlayListByShareCode", c.Request.URL.RawQuery, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 根据分享码拆分出游戏id
	dUid, pOpenId, shareType, gameId := service.IdShareCodeSplit(shareCode)
	if shareType != service.STGame {
		rsp.WriteMsg(model.CodeMsg{
			Code: model.ECBan,
			Msg:  "该邀请码不是游戏邀请码",
		})
		return
	}
	_, errCode, errMsg := service.ShareCodeCheck(dUid, pOpenId)
	if errCode != model.ErrCodeSuccess {
		rsp.WriteMsg(model.CodeMsg{Code: errCode, Msg: errMsg})
		return
	}
	// 查看该用户与这个游戏的对话记录
	dao.GetGameById(gameId)
	// 查看该游戏的当前状态 --是否被删除了,是否没有红包了，以及游戏的开场白等信息
	rsp.Data = &model.GamePlayInfo{
		GameName:    "",
		Prologue:    "",
		ChatList:    nil,
		GameState:   0,
		DigitalInfo: nil,
	}
	return
}

// 根据分享码与游戏聊天

// DelGame 删除游戏
func (l *LoginHandlerImp) DelGame(c *gin.Context) (rsp model.BaseRsp) {
	gameId := util.ToInt64(c.Query("gameId"))
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("DelGame", c.Request.URL.RawQuery, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 || gameId <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	//

	return
}
