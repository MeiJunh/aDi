package model

// FollowStatisticInfo 点赞关注统计信息
type FollowStatisticInfo struct {
	FollowNum int64 `json:"followNum"` // 关注数量
	FavorNum  int64 `json:"favorNum"`  // 点赞数量
	ViewNum   int64 `json:"viewNum"`   // 浏览数量
}

// FollowInfo 关注信息
type FollowInfo struct {
	UserBaseInfo
	FollowTime int64 `json:"followTime"` // 关注的时间
}

// FollowReq 关注、点赞入参
type FollowReq struct {
	Uid int64 `json:"uid"`
}

// DbFollow 关注信息
type DbFollow struct {
	Uid      int64 `json:"uid" db:"uid"`           // 被关注的人
	Follower int64 `json:"follower" db:"follower"` // 发起关注的人
	CTime    int64 `json:"ctime" db:"ctime"`       // 创建时间
}

// DbFavor 点赞信息
type DbFavor struct {
	Uid   int64 `json:"uid" db:"uid"`     // 被点赞的人
	Liker int64 `json:"liker" db:"liker"` // 发起点赞的人
	CTime int64 `json:"ctime" db:"ctime"` // 创建时间
}

// DbView 浏览信息
type DbView struct {
	Uid    int64 `json:"uid" db:"uid"`       // 被浏览的人
	Viewer int64 `json:"viewer" db:"viewer"` // 发起浏览的人
	CTime  int64 `json:"ctime" db:"ctime"`   // 创建时间
}

// DBFollowStatisticInfo 社交信息统计
type DBFollowStatisticInfo struct {
	Id        int64 `json:"id" db:"id"`
	Uid       int64 `json:"uid" db:"uid"`
	FollowNum int64 `json:"followNum" db:"follow_num"`
	FavorNum  int64 `json:"favorNum" db:"favor_num"`
	ViewNum   int64 `json:"viewNum" db:"view_num"`
}
