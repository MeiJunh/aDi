package dao

import (
	"aDi/log"
	"aDi/model"
	"github.com/jmoiron/sqlx"
	"math/rand"
)

// GetUserInfoByOpenID 根据open id获取用户信息
func GetUserInfoByOpenID(openID string) (user *model.UserInfo, err error) {
	user = &model.UserInfo{}
	err = dbClient.FindOneWithNull("select uid,nick,icon,age,sex,open_id,phone,session_key,union_id from t_user where open_id = ?", openID)
	if err != nil {
		log.Errorf("find user by open_id fail,err:%s", err.Error())
		return user, err
	}
	return user, err
}

// AddUserInfo 添加用户信息
func AddUserInfo(user *model.UserInfo) (err error) {
	err = Transaction(dbClient.DB, func(tx *sqlx.Tx) error {
		tmpUid := int64(0)
		errTmp := tx.Get(&tmpUid, "select uid from t_user order by uid desc limit 1")
		if errTmp != nil {
			log.Errorf("get max uid fail,err:%s", errTmp.Error())
			return errTmp
		}
		if tmpUid == 0 {
			// 初始值3000000000
			tmpUid = 3000000000
		}
		user.Uid = tmpUid + rand.Int63n(1000) // 1000内随机递增
		_, errTmp = tx.Exec("insert into t_user (uid,nick,icon,age,sex,open_id,union_id,phone,session_key) vaule (?,?,?,?,?,?,?,?)", user.Uid, user.Nick,
			user.Icon, user.Age, user.Sex, user.OpenID, user.UnionID, user.Phone, user.SessionKey)
		if errTmp != nil {
			log.Errorf("add user fail,err:%s", errTmp.Error())
			return errTmp
		}
		return nil
	})
	if err != nil {
		log.Errorf("add user fail,err:%s", err.Error())
		return err
	}
	return nil
}

// UpdateUserSessionKey 更新用户session key信息
func UpdateUserSessionKey(sessionKey string, uid int64) (err error) {
	_, err = dbClient.Exec("update user set session_key = ? where uid = ?", sessionKey, uid)
	if err != nil {
		log.Errorf("update user session key err:%v", err)
		return err
	}
	return nil
}
