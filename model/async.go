package model

// DbAsyncResult 异步结果
type DbAsyncResult struct {
	ID     int64  `json:"id" db:"id"`         // id
	UID    int64  `json:"uid" db:"uid"`       // 聊天人uid
	Result string `json:"result" db:"result"` // 结果
}
