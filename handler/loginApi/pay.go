package loginApi

import (
	"aDi/handler/comm"
	"aDi/model"
	"aDi/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

// QueryOrder 订单信息查询
func (l *LoginHandlerImp) QueryOrder(c *gin.Context) (rsp model.BaseRsp) {
	tradeNo := c.Query("tradeNo")
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("QueryOrder", tradeNo, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 查找订单信息 -- 如果成功或者失败则不在调用微信接口

	// 否则调用微信接口进行查询

	rsp.Data = &model.PayResult{
		TradeState:     0,
		TradeStateDesc: "",
	}
	return
}
