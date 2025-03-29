package service

import (
	"aDi/config"
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"context"
	"crypto/rsa"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/partnerpayments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"time"
)

// wxMchClient 微信直连商户client
var wxMchClient *core.Client
var MchPrivateKey *rsa.PrivateKey
var MchPublicKey *rsa.PublicKey

// https://pay.weixin.qq.com/doc/v3/partner/4012925289
func initWxMchClient() {
	var err error
	MchPrivateKey, err = utils.LoadPrivateKey(config.GetWxPayConf().WxDirectApiClientKey)
	if err != nil {
		log.Error("load direct merchant private key error!")
		panic("load direct merchant private key error!")
		return
	}

	MchPublicKey, err = utils.LoadPublicKey(config.GetWxPayConf().WxPublicKeyStr)
	if err != nil {
		panic(fmt.Errorf("load wechatpay public key err:%s", err.Error()))
	}
	// 加载商户私钥，商户私钥会用来生成请求的签名
	ctx := context.Background()
	// 使用商户私钥等初始化 wxMchClient，并使它具有自动定时获取微信支付平台证书的能力
	//opts := []core.ClientOption{
	//	option.WithWechatPayAutoAuthCipher(config.GetWxPayConf().WxDirectMchId, config.GetWxPayConf().WxDirectMchCertificateSerialNumber, MchPrivateKey, config.GetWxPayConf().WxDirectAPIv3Key),
	//}
	opts := []core.ClientOption{
		option.WithWechatPayPublicKeyAuthCipher(config.GetWxPayConf().WxDirectMchId, config.GetWxPayConf().WxDirectMchCertificateSerialNumber, MchPrivateKey, config.GetWxPayConf().WxPublicKeyId, MchPublicKey),
	}
	wxMchClient, err = core.NewClient(ctx, opts...)
	if err != nil {
		log.Errorf("new wx client failed! err: %s", err)
		panic("new wx pay client failed!")
	}
	return
}

// MchCertVisitorInstance 直连商户获取证书访问器
func MchCertVisitorInstance() (certVisitor core.CertificateVisitor, err error) {
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, MchPrivateKey, config.GetWxPayConf().WxDirectMchCertificateSerialNumber,
		config.GetWxPayConf().WxDirectMchId, config.GetWxPayConf().WxDirectAPIv3Key)
	if err != nil {
		log.Errorf("register cert downloader failed! err: %s", err.Error())
		return
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certVisitor = downloader.MgrInstance().GetCertificateVisitor(config.GetWxPayConf().WxDirectMchId)
	return
}

// WxMchJsapi jsapi下单
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_1.shtml
func WxMchJsapi(ctx context.Context, input *model.UnifiedOrderReq) (order *model.PayInfo, err error) {
	tradNo := GenOutTradeNo()
	svc := jsapi.JsapiApiService{Client: wxMchClient}
	appID := config.GetWxPayConf().AppId
	resp, httpRsp, err := svc.PrepayWithRequestPayment(ctx, jsapi.PrepayRequest{
		SpAppid:       core.String(appID),
		SpMchid:       core.String(config.GetWxPayConf().WxDirectMchId),
		Description:   core.String(input.ProductDescription),
		OutTradeNo:    core.String(tradNo),
		TimeExpire:    getWxOrderExpireTime("60"),
		NotifyUrl:     core.String(config.GetWxPayConf().WxDirectPayNotify),
		SupportFapiao: core.Bool(true),
		Amount: &jsapi.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(input.Amount),
		},
		Payer: &jsapi.Payer{
			SpOpenid: core.String(input.OpenId),
		},
	}, "")
	if err != nil {
		log.Errorf("wechat add js api fail,input:%+v,rsp:%+v,err:%s", input, httpRsp, err.Error())
		return order, err
	}
	order = &model.PayInfo{
		Appid:     appID,
		Timestamp: util.ToInt64(*resp.TimeStamp),
		Nonce:     *resp.NonceStr,
		Package:   *resp.Package,
		SignType:  "MD5",
		PaySign:   *resp.PaySign,
		PrepayId:  *resp.PrepayId,
		TradeNo:   tradNo,
	}
	// 添加订单信息
	err = dao.AddPayOrder(tradNo, input)
	if err != nil {
		log.Errorf("add js api fail,input:%+v,err:%s", input, err.Error())
		return order, err
	}
	return order, err
}

// GenOutTradeNo 生成的微信内部订单号
func GenOutTradeNo() (tradeNo string) {
	return util.GenOutTradeNo("WX")
}

func getWxOrderExpireTime(min string) *time.Time {
	d, _ := time.ParseDuration(min + "m")
	expireTime := time.Now().Add(d)
	return core.Time(expireTime)
}

// WxPayStateStringAdapter 微信支付状态转换
func WxPayStateStringAdapter(tradeState *string) (state model.PayState) {
	switch *tradeState {
	case "SUCCESS":
		state = model.PaySuccess
	case "REFUND", "NOTPAY", "CLOSED", "REVOKED", "USERPAYING", "PAYERROR":
		state = model.PayFail
	default:
		state = model.PayNotPay
	}
	return
}

// DealPay 处理所有的购买逻辑
func DealPay(order *model.DbTPayCenter) {
	if order.TradeState != model.PaySuccess {
		// 不是成功的不处理
		log.Debugf("order trade state no success,order:%s", util.MarshalToStringWithOutErr(order))
		return
	}
	switch order.ProdType {
	case model.PTChatNum:
		// 数字人聊天次数
		DealCharNumBuy(order)
	case model.PTRe:
		// 红包
		DealGameAddBuy(order)
	default:
		// 不在对应的业务内 -- 打印日志
		log.Errorf("no pay type deal,order info:%s", util.MarshalToStringWithOutErr(order))
	}
	return
}

// DealCharNumBuy 处理数字人对话购买
func DealCharNumBuy(order *model.DbTPayCenter) {
	// 进行扩展数据解析
	expandStruct := &model.BuyChatNumReq{}
	_ = jsoniter.UnmarshalFromString(order.ExpandStr, expandStruct)

	// 判断该用户
	effect, err := dao.UpsertConversation(&model.DbConversationConf{
		UID:            order.Uid,
		DigitalUID:     expandStruct.DigitalUid,
		IsAnonymity:    expandStruct.IsAnonymity,
		ChatConf:       util.MarshalToStringWithOutErr(expandStruct.ChatConfStruct),
		ChatTotalNum:   expandStruct.BuyNum * 10, // 目前一次购买是10次对话
		ChatConfStruct: expandStruct.ChatConfStruct,
	})
	if err != nil {
		log.Errorf("deal char num buy fail,err:%s", err.Error())
		return
	}
	log.Debugf("deal char num buy success,trade no:%s,effect：%d", order.TradeNo, effect)
	return
}

// DealGameAddBuy 处理游戏创建
func DealGameAddBuy(order *model.DbTPayCenter) {
	// 扩展参数解析
	expandStruct := &model.GameInfo{}
	_ = jsoniter.UnmarshalFromString(order.ExpandStr, expandStruct)
	// 进行游戏创建
	effect, err := dao.AddGame(order.Uid, expandStruct)
	log.Debugf("add game,trade no:%s,effect:%d,err:%v", order.TradeNo, effect, err)
	return
}
