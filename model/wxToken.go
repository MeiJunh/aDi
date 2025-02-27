package model

// WxTokenInfo 微信的access token信息
type WxTokenInfo struct {
	Token    string `json:"token" db:"token"`
	ExpireAt int64  `json:"expireAt" db:"expire_at"`
	IsLocked bool   `json:"isLocked" db:"is_locked"`
}

// GetUserPhoneNumberRsp 微信获取手机号返回结果
type GetUserPhoneNumberRsp struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}
