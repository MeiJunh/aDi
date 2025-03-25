package model

import jsoniter "github.com/json-iterator/go"

// DBUserInfo 用户信息
type DBUserInfo struct {
	Uid        int64  `json:"uid" db:"uid"`              // 用户id
	Nick       string `json:"nick" db:"nick"`            // 昵称
	Icon       string `json:"icon" db:"icon"`            // 头像
	Age        int64  `json:"age" db:"age"`              // 年龄
	Sex        string `json:"sex" db:"sex"`              // 性别
	Education  string `json:"education" db:"education"`  // 学历
	OpenID     string `json:"-" db:"open_id"`            // open id
	UnionID    string `json:"-" db:"union_id"`           // union id
	MBTIStr    string `json:"-" db:"mbti_str"`           // mkti信息
	TagInfoStr string `json:"-" db:"tag_info_str"`       // 标签信息
	Expand     string `json:"expand" db:"expand"`        // 拓展信息
	Phone      string `json:"-" db:"phone"`              // 电话
	SessionKey string `json:"-" db:"session_key"`        // 微信隐式登录返回的session key -- 先记录，暂时没有啥用
	IsVisible  Switch `json:"isVisible" db:"is_visible"` // 是否可见
}

// UnRegisterInfo 用户未注册信息
type UnRegisterInfo struct {
	Id         int64  `json:"id" db:"id"`
	OpenID     string `json:"-" db:"open_id"`
	SessionKey string `json:"-" db:"session_key"`
	UnionID    string `json:"-" db:"union_id"`
}

// UserBaseInfo 用户基础信息
type UserBaseInfo struct {
	Uid       int64  `json:"uid"`       // 用户id
	Nick      string `json:"nick"`      // 昵称
	Icon      string `json:"icon"`      // 头像
	Age       int64  `json:"age"`       // 年龄
	Sex       string `json:"sex"`       // 性别
	Education string `json:"education"` // 学历
	Visible   Switch `json:"visible"`
}

// UserAllInfo 用户的全部信息
type UserAllInfo struct {
	UserBaseInfo              // 用户基础信息
	MBTI         *MBTIInfo    `json:"mbti"`    // MBTI 信息
	TagInfo      *UserTagInfo `json:"tagInfo"` // 标签信息
	UserExpand                // 扩展信息
}

func (u *UserAllInfo) GetIcon() string {
	if u == nil {
		return ""
	}
	return u.Icon
}

// TransDbUserToAll 将数据库结构体进行转化all user
func TransDbUserToAll(info *DBUserInfo) (userAllInfo *UserAllInfo) {
	if info == nil {
		return &UserAllInfo{}
	}
	userAllInfo = &UserAllInfo{
		UserBaseInfo: *(TransDbUserToBase(info)),
		MBTI:         &MBTIInfo{},
		TagInfo:      &UserTagInfo{},
		UserExpand:   UserExpand{},
	}
	if info.MBTIStr != "" {
		_ = jsoniter.UnmarshalFromString(info.MBTIStr, userAllInfo.MBTI)
	}
	if info.TagInfoStr != "" {
		_ = jsoniter.UnmarshalFromString(info.TagInfoStr, userAllInfo.TagInfo)
	}
	if info.Expand != "" {
		_ = jsoniter.UnmarshalFromString(info.Expand, &userAllInfo.UserExpand)
	}
	return userAllInfo
}

// TransDbUserToBase 将数据库结构体转化成base
func TransDbUserToBase(info *DBUserInfo) (userBase *UserBaseInfo) {
	if info == nil {
		return &UserBaseInfo{}
	}
	return &UserBaseInfo{
		Uid:       info.Uid,
		Nick:      info.Nick,
		Icon:      info.Icon,
		Age:       info.Age,
		Sex:       info.Sex,
		Education: info.Education,
		Visible:   info.IsVisible,
	}
}

// UserExpand 用户扩展信息
type UserExpand struct {
	Location      string `json:"location"`      // 所在地 -- 省市用/连接
	HomeTown      string `json:"homeTown"`      // 家乡
	BirthDay      string `json:"birthDay"`      // 生日
	Constellation string `json:"constellation"` // 星座
	Job           string `json:"job"`           // 职业
	Company       string `json:"company"`       // 公司
	School        string `json:"school"`        // 学校
	Major         string `json:"major"`         // 专业
}

// MBTIInfo 信息
type MBTIInfo struct {
	One   string `json:"one"`
	Two   string `json:"two"`
	Three string `json:"three"`
	Four  string `json:"four"`
}

// 用户标签信息
type UserTagInfo struct {
	Hobby           []string `json:"hobby"`           // 爱好
	Skill           []string `json:"skill"`           // 技能
	Idol            []string `json:"idol"`            // 偶像
	Travel          []string `json:"travel"`          // 旅行经历
	IndividualTrait []string `json:"individualTrait"` // 个人特质
	Movie           []string `json:"movie"`           // 影视
	Music           []string `json:"music"`           // 音乐
	Book            []string `json:"book"`            // 阅读
	Sport           []string `json:"sport"`           // 运动
	Food            []string `json:"food"`            // 美食
}

type ModVisibleReq struct {
	Visible Switch `json:"visible"` // 是否可见
}
