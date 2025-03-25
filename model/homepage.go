package model

// DbHomepageConf 主页配置表
type DbHomepageConf struct {
	ID      int64   `json:"id" db:"id"`            //
	UID     int64   `json:"uid" db:"uid"`          // uid
	ConfStr string  `json:"confStr" db:"conf_str"` // 配置信息
	HpConf  *HpConf `json:"-"`                     // 配置信息
}

// Point 点信息
type Point struct {
	X float64 `json:"x"` // x轴
	Y float64 `json:"y"` // y轴
}

// HpConf 个人主页面配置信息
type HpConf struct {
	CardList []*HpCard `json:"cardList"` // 卡片信息
}

// HpCardType 个人页面卡片类型
type HpCardType int32

// HpCard 个人主页面的卡片信息
type HpCard struct {
	LeftTopPosition     Point      `json:"leftTopPosition"`     // 左上角的点信息
	RightBottomPosition Point      `json:"rightBottomPosition"` // 右下角配置信息
	Info                string     `json:"info"`                // 配置信息
	CardType            HpCardType `json:"cardType"`            // 卡片类型
}

// GetOtherHomePageInfoByShareCodeRsp 根据分享码获取其他人的主页信息 -- 包括主页配置与个人信息
type GetOtherHomePageInfoByShareCodeRsp struct {
	HpConfInfo    *HpConf              `json:"hpConfInfo"`    // 主页配置信息
	StatisticInfo *FollowStatisticInfo `json:"statisticInfo"` // 统计信息
	DigitalInfo   *DigitalShadowInfo   `json:"digitalInfo"`   // 数字人相关信息
}
