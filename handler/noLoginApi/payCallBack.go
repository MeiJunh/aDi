package noLoginApi

import (
	"aDi/config"
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

const (
	WxRespCodeSuccess = "SUCCESS"
	WxRespCodeError   = "ERROR"
)

type WxNotifyResponse struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

// MchWxPayNotify 直连商户支付回调通知接口
func MchWxPayNotify(c *gin.Context) {
	// 默认失败
	resp := WxNotifyResponse{Code: WxRespCodeError, Msg: "系统错误"}
	defer func() {
		c.JSON(200, resp)
	}()

	// 获取直连商户证书访问管理器
	certVisitor, err := service.MchCertVisitorInstance()
	if err != nil {
		log.Errorf("wx direct mch pay notify get cert visitor failed! err: %s")
		return
	}
	handler, err := notify.NewRSANotifyHandler(config.GetWxPayConf().WxDirectAPIv3Key, verifiers.NewSHA256WithRSACombinedVerifier(certVisitor, config.GetWxPayConf().WxPublicKeyId, *service.MchPublicKey))
	if err != nil {
		log.Errorf("wx direct mch pay notify get handler failed! err: %s", err.Error())
		return
	}
	// 处理回调逻辑，将支付回调通知中的内容，解析为 payments.Transaction。
	notifyReq := new(notify.Request)
	transaction := new(payments.Transaction)
	notifyReq, err = handler.ParseNotifyRequest(context.Background(), c.Request, transaction)
	if err != nil {
		log.Errorf("wx direct mch pay notify parse data failed!buff:%s, err: %s", string(util.CopyGetRequestBody(c.Request)), err)
		resp.Msg = "解析请求Body失败"
		return
	}
	log.Infof("parse wx direct mch notify data success! notifyReq: %+v, parseData: %+v", notifyReq, transaction)
	// 查出订单是否存在
	order, err := dao.GetPayOrderByTradeNo(*transaction.OutTradeNo)
	if err != nil {
		log.Errorf("wx direct mch pay notify get order failed! tradeNo: %s, err: %s", *transaction.OutTradeNo, err)
		resp.Msg = "获取订单信息失败"
		return
	}
	if order == nil || order.Id <= 0 {
		resp.Msg = "查询不到对应的订单"
		log.Errorf("get no order, tradeNo: %s", *transaction.OutTradeNo)
		return
	}
	// 订单已经是成功状态 -- 直接返回
	if order.TradeState == model.PaySuccess {
		resp.Code = WxRespCodeSuccess
		resp.Msg = "成功"
		return
	}

	// 更新订单状态，触发回调
	order.ChTradeNo = *transaction.TransactionId
	order.TradeState = service.WxPayStateStringAdapter(transaction.TradeState)
	order.StatusMsg = *transaction.TradeStateDesc
	order.PayerTotal = *transaction.Amount.PayerTotal
	t, _ := util.ParseTimerFromStr(*transaction.SuccessTime, util.TimestampFormatT)
	order.OrderEtime = t.Format(util.TimestampFormat)
	// 只有成功的订单会回调，这里不再判断状态
	effect, err := dao.UpdatePayOrder(order)
	if err != nil {
		resp.Msg = "更新订单状态出错"
		log.Errorf("wx direct mch notify update order info failed! tradeNo: %s", order.TradeNo)
		return
	}
	if effect <= 0 {
		// 订单没有发生变化直接返回
		resp.Code = WxRespCodeSuccess
		resp.Msg = "成功"
		return
	}
	log.Infof("wx direct mch notify update order info success! tradeNo: %s", order.TradeNo)
	// 执行业务逻辑
	util.GoSafe(func() {
		service.DealPay(order)
	})
	resp.Code = WxRespCodeSuccess
	resp.Msg = "成功"
	return
}
