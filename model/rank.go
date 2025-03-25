package model

// DbChatStatistic 数据库统计信息
type DbChatStatistic struct {
	DigitalUid  int64  `json:"digitalUid" db:"digital_uid"`   // 机器人对应的uid
	ChatUid     int64  `json:"chatUid" db:"chat_uid"`         // 发起聊天对应的uid
	IsAnonymity Switch `json:"isAnonymity" db:"is_anonymity"` // 是否是匿名
	ChatNum     int64  `json:"chatNum" db:"chat_num"`         // 聊天次数
}

// RankInfo 榜单信息
type RankInfo struct {
	UserBaseInfo
	Score int64 `json:"score"` // 榜单得分
}
