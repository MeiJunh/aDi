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
	// 根据分享码拆分出游戏id并且获取游戏信息
	_, rsp.Data, rsp.Code, rsp.Msg = service.GetGameInfoByShareCode(uid, shareCode)
	return
}

// ChatWithGame 根据分享码与游戏聊天
func (l *LoginHandlerImp) ChatWithGame(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.ChatWithGameReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ChatWithGame", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	// 参数解析
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.ShareCode == "" || req.Input == "" {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 根据分享码获取对应的游戏信息
	gameInfo, gamePlayInfo, errCode, errMsg := service.GetGameInfoByShareCode(uid, req.ShareCode)
	if errCode != model.ErrCodeSuccess {
		rsp.WriteCodeMsg(errCode, errMsg)
		log.Errorf("get game info fail,err:%s", err.Error())
		return
	}
	// 如果当前游戏的状态不是正常 -- 则直接返回
	if gamePlayInfo.GameState != model.GSDefault {
		errMsg = "当前游戏已经没有奖励，请更换对应的游戏邀请码重试"
		switch gamePlayInfo.GameState {
		case model.GSDel: // 被删除
			errMsg = "当前游戏已被删除"
		case model.GSNoRE: // 红包已经发完
			errMsg = "当前游戏红包已经发完"
		case model.GSExpire: // 游戏过期
			errMsg = "当前游戏已过期"
		case model.GSPlayOver: // 当前用户已经将该游戏玩完了
			errMsg = "您当前所在游戏已经达到游戏次数上限"
		case model.GSPlayWin: // 当前用户已经获得了该游戏的奖励
			errMsg = "您已经获得了该游戏的奖励"
		}
		rsp.WriteCodeMsg(model.ECBan, errMsg)
		return
	}
	// 添加游戏信息
	errCode, errMsg = service.ChatWithGame(uid, req.Input, gameInfo)
	rsp.WriteCodeMsg(errCode, errMsg)
	return
}

// DelGame 删除游戏
func (l *LoginHandlerImp) DelGame(c *gin.Context) (rsp model.BaseRsp) {
	gameId := util.ToInt64(c.Query("gameId"))
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("DelGame", c.Request.URL.RawQuery, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 || gameId <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取游戏信息
	gameInfo, err := dao.GetGameById(gameId)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get game info fail,err:%s", err.Error())
		return
	}
	if gameInfo == nil || gameInfo.Id <= 0 || gameInfo.State == model.GSDel {
		// 游戏不存在或者游戏已经被删除则直接返回
		return
	}
	if gameInfo.Uid != uid {
		// 如果游戏不属于当前用户 -- 直接报错
		rsp.WriteMsg(model.CodeMsg{Code: model.ECBan, Msg: "游戏不属于当前用户，请刷新重选"})
		return
	}
	// 删除游戏
	errCode, errMsg := service.DelGame(gameInfo)
	if errCode != model.ErrCodeSuccess {
		rsp.WriteMsg(model.CodeMsg{Code: errCode, Msg: errMsg})
		log.Errorf("del game info fail,err code:%d,err msg:%s", errCode, errMsg)
		return
	}
	return
}
