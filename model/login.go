package model

// LoginReq 登录参数
type LoginReq struct {
	Code string `json:"code" binding:"required"`
}

// LoginRsp 登录回参
type LoginRsp struct {
	UnRegisterId int64      `json:"unRegisterId"` // 如果用户为未注册时，使用该id
	ApiKey       string     `json:"apiKey"`       // 用于参数校验的api key
	User         *UserInfo  `json:"user"`         // 用户信息
	TokenInfo    *TokenInfo `json:"tokenInfo"`    // token信息
}

type TokenInfo struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expireAt"`
}

// RegisterReq 注册入参
type RegisterReq struct {
	UnRegisterId  int64  `json:"unRegisterId"`  // 未注册用户id
	WxDynamicCode string `json:"wxDynamicCode"` // 微信动态令牌
}

// RegisterRsp 注册回参
type RegisterRsp struct {
	User      *UserInfo  `json:"user"`      // 用户信息
	TokenInfo *TokenInfo `json:"tokenInfo"` // token信息
}
