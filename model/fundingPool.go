package model

// FundingPoolType 资金池类型
type FundingPoolType int32

// TradeType 交易类型
type TradeType int32

const (
	FPTDigitalInput = FundingPoolType(1) // 数字人收入
	FPTRe           = FundingPoolType(2) // 红包资金池

	TTInCome  = TradeType(1) // 进账
	TTOutCome = TradeType(2) // 出账
)

// FundingPoolInfo 资金池信息
type FundingPoolInfo struct {
	PoolType    FundingPoolType `json:"poolType"`    // 资金池类型1：数字人收入，2：红包
	TotalAmount int64           `json:"totalAmount"` // 剩余总金额
}

// WithdrawFundingPoolReq 资金池提现参数
type WithdrawFundingPoolReq struct {
	PoolType       FundingPoolType `json:"poolType"`       // 资金池类型1：数字人收入，2：红包
	WithdrawAmount int64           `json:"withdrawAmount"` // 提现金额
}

// FPDetailInfo 资金池详细信息
type FPDetailInfo struct {
	TradeNo   string    `json:"tradeNo"`   // 交易单号
	Amount    int64     `json:"amount"`    // 金额--单位分,为正数
	TradeType TradeType `json:"tradeType"` // 交易类型，判断金额前面是否需要加符号 1进账，2出账
	TradeTime int64     `json:"TradeTime"` // 时间--秒级时间戳
}

// DbFundingPool 资金池数据表
type DbFundingPool struct {
	ID                int64           `json:"id" db:"id"`                                 //
	Uid               int64           `json:"uid" db:"uid"`                               // uid
	PoolType          FundingPoolType `json:"poolType" db:"pool_type"`                    // 资金池类型:1:数字人,2:红包
	AllTotalAmount    int64           `json:"allTotalAmount" db:"all_total_amount"`       // 一直以来的总收入 -- 单位分
	AllWithdrawAmount int64           `json:"allWithdrawAmount" db:"all_withdraw_amount"` // 一直以来被提现的金额 -- 单位分
	CurTotalAmount    int64           `json:"curTotalAmount" db:"cur_total_amount"`       // 当前资金 -- 单位分
	CreateTime        int64           `json:"createTime" db:"create_time"`                //
}

// DbFundingPoolDetail 资金池记详情表
type DbFundingPoolDetail struct {
	ID          int64           `json:"id" db:"id"`                     //
	Uid         int64           `json:"uid" db:"uid"`                   // uid
	PoolType    FundingPoolType `json:"poolType" db:"pool_type"`        // 资金池类型:1:数字人,2:红包
	TradeNo     string          `json:"tradeNo" db:"trade_no"`          // 业务订单号--虚拟的
	TradeType   TradeType       `json:"tradeType" db:"trade_type"`      // 交易类型,1:进账,2:出账
	RealTradeNo string          `json:"realTradeNo" db:"real_trade_no"` // 业务订单号--真实的,用户看不到
	Amount      int64           `json:"amount" db:"amount"`             // 金额 -- 单位分
	CreateTime  int64           `json:"createTime" db:"create_time"`    //
}
