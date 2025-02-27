package model

// UserInfo 用户信息
type UserInfo struct {
	Uid        int64  `json:"uid" db:"uid"`       // 用户id
	Nick       string `json:"nick" db:"nick"`     // 昵称
	Icon       string `json:"icon" db:"icon"`     // 头像
	Age        int64  `json:"age" db:"age"`       // 年龄
	Sex        string `json:"sex" db:"sex"`       // 性别
	OpenID     string `json:"-" db:"open_id"`     // open id
	UnionID    string `json:"-" db:"union_id"`    // union id
	Phone      string `json:"-" db:"phone"`       // 电话
	SessionKey string `json:"-" db:"session_key"` // 微信隐式登录返回的session key -- 先记录，暂时没有啥用
}

// UnRegisterInfo 用户未注册信息
type UnRegisterInfo struct {
	Id         int64  `json:"id" db:"id"`
	OpenID     string `json:"-" db:"open_id"`
	SessionKey string `json:"-" db:"session_key"`
	UnionID    string `json:"-" db:"union_id"`
}
