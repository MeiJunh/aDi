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

// GetShareCode 获取分享码
func (l *LoginHandlerImp) GetShareCode(c *gin.Context) (rsp model.BaseRsp) {
	uid := comm.GetUidFromCon(c)
	shareType := util.ToInt64(c.Query("shareType"))
	id := util.ToInt64(c.Query("id"))
	defer util.TimeCost("GetShareCode", c.Request.URL.RawQuery, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取用户信息
	userInfo, err := dao.GetUserInfoByUid(uid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get user info fail,err:%s", err.Error())
		return
	}
	if userInfo == nil || userInfo.Uid <= 0 {
		rsp.WriteMsg(model.CodeMsg{Code: model.ECNoExist, Msg: "用户信息不存在"})
		return
	}
	var shareCode string
	switch service.ShareType(shareType) {
	case service.STGame:
		// 分享游戏 -- 需要判断这个游戏是不是该用户的 -- 如果是的话在构建分享码
		gameInfo, err := dao.GetGameById(id)
		if err != nil {
			rsp.WriteMsg(model.ErrIDbFail)
			log.Errorf("get game info fail,err:%s", err.Error())
			return
		}
		if gameInfo == nil || gameInfo.Id <= 0 {
			rsp.WriteMsg(model.CodeMsg{Code: model.ECNoExist, Msg: "该游戏不存在"})
			return
		}
		if gameInfo.Uid != uid {
			rsp.WriteMsg(model.CodeMsg{Code: model.ECBan, Msg: "该游戏不属于当前用户，请分享自己的游戏"})
			return
		}
		// 获取带id的分享码
		shareCode = service.GetShareCodeByUidAndId(uid, userInfo.OpenID, service.ShareType(shareType), id)
	default:
		shareCode = service.GetShareCodeByUid(uid, userInfo.OpenID, service.ShareType(shareType))
	}
	// 获取分享码
	rsp.Data = shareCode
	return
}
