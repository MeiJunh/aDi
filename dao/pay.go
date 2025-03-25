package dao

import (
	"aDi/log"
	"aDi/model"
	"time"
)

// AddPayOrder 添加订单
func AddPayOrder(tradeNo string, input *model.UnifiedOrderReq) (err error) {
	_, err = dbClient.Exec("INSERT INTO t_pay_center (uid, open_id, trade_no, amount,prod_type, prod_desc,"+
		" prod_attach, order_ctime, trade_state, expand_str) VALUES "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", input.Uid, input.OpenId, tradeNo, input.Amount, input.ProdType,
		input.ProductDescription, input.ProductAttach, time.Now().Format(time.DateTime), model.PayNotPay, input.ExpandStr)
	if err != nil {
		log.Errorf("add pay order fail,trade no:%s,input:%v,err:%s", tradeNo, input, err.Error())
		return err
	}
	return err
}

// GetPayOrderByTradeNo 根据订单id获取订单信息
func GetPayOrderByTradeNo(tradeNo string) (order *model.DbTPayCenter, err error) {
	order = &model.DbTPayCenter{}
	query := "SELECT id, uid, open_id, trade_no, ch_trade_no, amount, payer_total, prod_type, prod_desc, prod_attach, order_ctime, order_life_time, " +
		"order_etime, trade_state, expand_str, status_msg, UNIX_TIMESTAMP(create_time) AS create_time, UNIX_TIMESTAMP(update_time) AS update_time FROM t_pay_center " +
		"where trade_no = ?"
	err = dbClient.FindOneWithNull(order, query, tradeNo)
	if err != nil {
		log.Errorf("get pay order fail,tradeNo:%s,err:%s", tradeNo, err)
		return order, err
	}
	return order, err
}

// UpdatePayOrder 更新订单
func UpdatePayOrder(order *model.DbTPayCenter) (effect int64, err error) {
	query := "update t_pay_center set ch_trade_no = ?,payer_total = ?,trade_state = ?,status_msg = ? where trade_no = ?"
	result, err := dbClient.Exec(query, order.ChTradeNo, order.PayerTotal, order.TradeState, order.StatusMsg, order.TradeNo)
	if err != nil {
		log.Errorf("update pay order fail,tradeNo:%s,err:%s", order.TradeNo, err)
		return effect, err
	}
	effect, err = result.RowsAffected()
	if err != nil {
		log.Errorf("update pay order fail,tradeNo:%s,err:%s", order.TradeNo, err)
		return effect, err
	}
	return effect, err
}
