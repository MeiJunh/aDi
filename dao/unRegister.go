package dao

import (
	"aDi/log"
	"aDi/model"
)

// UpsertUnRegisterInfo 添加未注册信息 -- 主要是插入open id和用session key的关系,用户注册时根据open id获取session key
func UpsertUnRegisterInfo(openId, unionId, sessionKey string) (insertId int64, err error) {
	// 先查找未注册用户信息，如果存在则更新session key
	unRegisterInfo, err := GetUnregisterInfoByOpenId(openId)
	if err != nil {
		log.Errorf("get unregister info fail,err:%s", err.Error())
		return insertId, err
	}
	if unRegisterInfo != nil && unRegisterInfo.Id > 0 {
		// 更新未注册用户的session key
		err = UpdateUnRegisterSessionKey(unRegisterInfo.Id, sessionKey)
		if err != nil {
			log.Errorf("update unregister session key fail,err:%s", err.Error())
			return insertId, err
		}
		return unRegisterInfo.Id, nil
	}
	// 未注册用户信息插入
	insertId, err = AddUnRegisterInfo(unionId, unionId, sessionKey)
	if err != nil {
		log.Errorf("add unregister info fail,err:%s", err.Error())
		return insertId, err
	}
	return insertId, err
}

// GetUnregisterInfoByOpenId 根据open id获取未注册用户信息
func GetUnregisterInfoByOpenId(openId string) (unRegisterInfo *model.UnRegisterInfo, err error) {
	return unRegisterInfo, err
}

// UpdateUnRegisterSessionKey 更新未注册用户信息的session key
func UpdateUnRegisterSessionKey(id int64, sessionKey string) (err error) {
	return err
}

// AddUnRegisterInfo 添加未注册用户信息
func AddUnRegisterInfo(openId, unionId, sessionKey string) (insertId int64, err error) {
	return insertId, err
}
