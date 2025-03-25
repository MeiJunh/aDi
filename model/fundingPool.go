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
