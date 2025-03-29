package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"github.com/jmoiron/sqlx"
)

// FundingPoolAddAmount 添加资金池金额 -- 进账
func FundingPoolAddAmount(detail *model.DbFundingPoolDetail) (err error) {
	err = Transaction(dbClient.DB, func(tx *sqlx.Tx) error {
		// 进行资金池总金额修改
		errT := fundingPoolAddAmount(tx, detail.Uid, detail.PoolType, detail.Amount)
		if errT != nil {
			log.Errorf("funding pool add amount err: %s", errT.Error())
			return errT
		}
		// 进行资金池详情信息添加
		errT = fundingPoolAddDetail(tx, detail)
		if errT != nil {
			log.Errorf("funding pool add detail err: %s", errT.Error())
			return errT
		}
		return nil
	})
	if err != nil {
		log.Errorf("funding pool add amount err: %s", err.Error())
		return err
	}
	return nil
}

// fundingPoolAddAmount 资金池新增收入
func fundingPoolAddAmount(tx *sqlx.Tx, uid int64, poolType model.FundingPoolType, amount int64) (err error) {
	_, err = tx.Exec("INSERT INTO t_funding_pool (uid, pool_type, all_total_amount, cur_total_amount) VALUES (?, ?, ?, ?) "+
		"on duplicate key update all_total_amount = all_total_amount + values(all_total_amount),cur_total_amount = cur_total_amount + values(cur_total_amount)",
		uid, poolType, amount, amount)
	if err != nil {
		log.Errorf("funding pool add amount fail,err:%s", err.Error())
		return err
	}
	return err
}

// fundingPoolAddDetail 添加资金池明细
func fundingPoolAddDetail(tx *sqlx.Tx, detail *model.DbFundingPoolDetail) (err error) {
	if detail.TradeNo == "" {
		// 如果没有订单 -- 则初始化一个
		detail.TradeNo = util.GenerateFundingPoolTradeNo(detail.PoolType)
	}
	_, err = tx.Exec("INSERT INTO t_funding_pool_detail (uid, pool_type, trade_no, trade_type, real_trade_no, amount) VALUES (?, ?, ?, ?, ?, ?)",
		detail.Uid, detail.PoolType, detail.TradeNo, detail.TradeType, detail.RealTradeNo, detail.Amount)
	if err != nil {
		log.Errorf("funding pool detail fail,err:%s", err.Error())
		return err
	}
	return err
}
