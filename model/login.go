package model

// LoginReq 登录参数
type LoginReq struct {
	Code string `json:"code" binding:"required"`
}

// LoginRsp 登录回参
type LoginRsp struct {
	UnRegisterId int64      `json:"unRegisterId"` // 如果用户为未注册时，使用该id
	KeyCode      string     `json:"keyCode"`      // 用于参数校验的keyCode
	User         *UserInfo  `json:"user"`         // 用户信息
	TokenInfo    *TokenInfo `json:"tokenInfo"`    // token信息
}

type TokenInfo struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire"`
}
