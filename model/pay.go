package model

// PayState 支付状态
type PayState int32

const (
	PayNotPay  = 0 // 未支付
	PaySuccess = 1 // 支付成功
	PayFail    = 2 // 支付失败
)

// PayInfo 下订单之后的信息
type PayInfo struct {
	Appid     string `json:"appid,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Package   string `json:"package,omitempty"`
	SignType  string `json:"signType,omitempty"`
	PaySign   string `json:"paySign,omitempty"`
	PrepayId  string `json:"prepayId,omitempty"`
	TradeNo   string `json:"tradeNo,omitempty"`
}

// PayResult 支付结果
type PayResult struct {
	TradeState     PayState `json:"tradeState"`
	TradeStateDesc string   `json:"tradeStateDesc"`
}

// UnifiedOrderReq 下单需要的参数
type UnifiedOrderReq struct {
	Uid                int64    `json:"uid"`                // uid
	OpenId             string   `json:"openId"`             // 用户open id
	ProductDescription string   `json:"productDescription"` // 商品描述
	ProductAttach      string   `json:"productAttach"`      // 商品附件
	Amount             int64    `json:"amount"`             // 收款金额 -- 单位分
	ProdType           ProdType `json:"prodType"`           // 商品类型
	ExpandStr          string   `json:"expandStr"`          // 拓展信息 -- 根据商品类型进行不同的解析
}

// ProdType 商品类型
type ProdType int32

const (
	PTChatNum = ProdType(1) // 聊天次数
	PTRe      = ProdType(2) // 红包类型
)

type DbTPayCenter struct {
	Id            int64    `json:"id" db:"id"`                         // 主键
	Uid           int64    `json:"uid" db:"uid"`                       // uid
	OpenId        string   `json:"openId" db:"open_id"`                // open id
	TradeNo       string   `json:"tradeNo" db:"trade_no"`              // 业务订单号
	ChTradeNo     string   `json:"chTradeNo" db:"ch_trade_no"`         // 渠道订单id
	Amount        int64    `json:"amount" db:"amount"`                 // 支付金额
	PayerTotal    int64    `json:"payerTotal" db:"payer_total"`        // 用户支付的总金额,分
	ProdType      ProdType `json:"prodType" db:"prod_type"`            // 商品类型
	ProdDesc      string   `json:"prodDesc" db:"prod_desc"`            // 订单描述
	ProdAttach    string   `json:"prodAttach" db:"prod_attach"`        // 订单Attach
	OrderCtime    string   `json:"orderCtime" db:"order_ctime"`        // 支付请求时间
	OrderLifeTime int64    `json:"orderLifeTime" db:"order_life_time"` // 订单有效时间, 秒
	OrderEtime    string   `json:"orderEtime" db:"order_etime"`        // 支付完成时间
	TradeState    PayState `json:"tradeState" db:"trade_state"`        // 交易状态 PayState
	ExpandStr     string   `json:"expandStr" db:"expand_str"`          // 具体业务产生订单时附带的信息
	StatusMsg     string   `json:"statusMsg" db:"status_msg"`          // 状态描述
	CreateTime    int64    `json:"createTime" db:"create_time"`        // 创建时间
	UpdateTime    int64    `json:"updateTime" db:"update_time"`        // 更新时间
}
