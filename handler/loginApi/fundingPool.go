package loginApi

import (
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetMyFundingPoolInfo 查看资金池信息
func (l *LoginHandlerImp) GetMyFundingPoolInfo(c *gin.Context) (rsp model.BaseRsp) {
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyFundingPoolInfo", uid, &rsp.Code)()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取资金池信息
	rsp.Data = &model.FundingPoolInfo{
		PoolType:    0,
		TotalAmount: 0,
	}
	return
}

// GetFundingPoolDetail 获取指定资金池的明细详情
func (l *LoginHandlerImp) GetFundingPoolDetail(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetFundingPoolDetail", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	rsp.Data = &model.GetListRsp{
		List: make([]*model.FPDetailInfo, 0),
	}
	return
}

// WithdrawFundingPool 资金池金额提现
func (l *LoginHandlerImp) WithdrawFundingPool(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.WithdrawFundingPoolReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("WithdrawFundingPool", req, &rsp.Code, strconv.FormatInt(uid, 10))()
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
	return
}
