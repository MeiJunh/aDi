package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"time"
)

const (
	// WxTokenId 微信token存的id固定为1
	WxTokenId = 1
)

// GetTokenFromDB 获取微信token信息
func GetTokenFromDB() (tokenInfo *model.WxTokenInfo, err error) {
	tokenInfo = &model.WxTokenInfo{}
	err = dbClient.FindOneWithNull(tokenInfo, "SELECT token,expire_at,is_locked FROM wx_token where id = ?", WxTokenId)
	if err != nil {
		log.Errorf("get token from db fail,err:%s", err.Error())
		return nil, err
	}
	return tokenInfo, nil
}

// WxTokenTryLock 微信token更新抢锁
func WxTokenTryLock() bool {
	// 抢锁更新 -- 锁未被抢或者是过期时间已经只剩下五分钟不到,则直接获取锁
	result, err := dbClient.Exec("UPDATE wx_token SET is_locked = ? WHERE (is_locked = ? OR expire_at < ?) and id = ?",
		model.SwitchOn, model.SwitchOff, time.Now().Unix()+util.MinuteSecond*5, WxTokenId)
	if err != nil {
		log.Errorf("get token from db fail,err:%s", err.Error())
		return false
	}

	// 判断是否更新成功 -- 更新成功表示获取锁成功
	affected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("get token from db fail,err:%s", err.Error())
		return false
	}
	return affected > 0
}

// WxTokenUnlock 微信token释放锁 -- 尽力释放就可以
func WxTokenUnlock() {
	_, err := dbClient.Exec("UPDATE wx_token SET is_locked = ? where id = ?", model.SwitchOff, WxTokenId)
	if err != nil {
		log.Errorf("unLock token fail,err:%s", err.Error())
		return
	}
	return
}

// UpdateWxToken 更新微信token信息
func UpdateWxToken(token string, expireAt int64) error {
	_, err := dbClient.Exec("INSERT INTO wx_token (id,token, expire_at, is_locked) VALUES (?,?,?,?) on duplicate key update token = values(token),"+
		"expire_at = values(expire_at),is_locked = values(is_locked)", WxTokenId, token, expireAt, model.SwitchOff)
	if err != nil {
		log.Errorf("update token fail,err:%s", err.Error())
		return err
	}
	return err
}
