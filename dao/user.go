package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math/rand"
)

// GetUserInfoByOpenID 根据open id获取用户信息
func GetUserInfoByOpenID(openID string) (user *model.DBUserInfo, err error) {
	user, err = getUserInfo(" where open_id = ?", openID)
	if err != nil {
		log.Errorf("find user by open_id fail,open id:%s,err:%s", openID, err.Error())
		return user, err
	}
	return user, err
}

// GetUserInfoByUid 根据uid获取用户信息
func GetUserInfoByUid(uid int64) (user *model.DBUserInfo, err error) {
	user, err = getUserInfo(" where uid = ?", uid)
	if err != nil {
		log.Errorf("find user by uid fail,uid:%d,err:%s", uid, err.Error())
		return user, err
	}
	return user, err
}

// AddUserInfo 添加用户信息
func AddUserInfo(user *model.DBUserInfo) (err error) {
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
		_, errTmp = tx.Exec("insert into t_user (uid,nick,icon,age,sex,open_id,union_id,phone,session_key,education,mbti_str,tag_info_str,expand) "+
			"vaule (?,?,?,?,?,?,?,?,'','','','')", user.Uid, user.Nick, user.Icon, user.Age, user.Sex, user.OpenID,
			user.UnionID, user.Phone, user.SessionKey)
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
	// 创建uid时进行该用户统计数据初始化
	InitSocialStatisticInfo(user.Uid)
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

// getUserInfo 获取用户信息
func getUserInfo(whereQuery string, params ...interface{}) (user *model.DBUserInfo, err error) {
	user = &model.DBUserInfo{}
	err = dbClient.FindOneWithNull(user, "select uid,nick,icon,age,sex,open_id,phone,session_key,union_id,education,mbti_str,"+
		"tag_info_str,expand,is_visible from t_user "+whereQuery, params...)
	if err != nil {
		log.Errorf("find user fail,err:%s", err.Error())
		return user, err
	}
	return user, err
}

// GetBUserListById 根据uid列表获取用户基础信息
func GetBUserListById(uidList []int64) (userList []*model.UserBaseInfo, err error) {
	if len(uidList) == 0 {
		return userList, err
	}
	userList, err = getBaseUserList(fmt.Sprint(" where uid in (%s)", util.JoinInt64(uidList, ",")))
	if err != nil {
		log.Errorf("get user fail,err:%s", err.Error())
		return userList, err
	}
	return userList, err
}

// getBaseUserList 获取用户基础信息
func getBaseUserList(whereQuery string, params ...interface{}) (userList []*model.UserBaseInfo, err error) {
	list := make([]*model.DBUserInfo, 0)
	err = dbClient.FindOneWithNull(&list, "select uid,nick,icon,age,sex,education,is_visible from t_user "+whereQuery, params...)
	if err != nil {
		log.Errorf("find user fail,err:%s", err.Error())
		return userList, err
	}
	for i := range list {
		userList = append(userList, &model.UserBaseInfo{
			Uid:       list[i].Uid,
			Nick:      list[i].Nick,
			Icon:      list[i].Icon,
			Age:       list[i].Age,
			Sex:       list[i].Sex,
			Education: list[i].Education,
			Visible:   list[i].IsVisible,
		})
	}
	return userList, err
}

// UpdateUserBaseInfo 更新用户基础信息 -- 除了标签和mbti信息不修改外，其余都进行修改
func UpdateUserBaseInfo(user *model.UserAllInfo) (err error) {
	_, err = dbClient.Exec("update t_user set nick = ?,icon = ?,age = ?,sex = ?,education = ?,expand = ? where uid = ?",
		user.Nick, user.Icon, user.Age, user.Sex, user.Education, util.MarshalToStringWithOutErr(user.UserExpand), user.Uid)
	if err != nil {
		log.Errorf("update user fail,err:%s", err.Error())
		return err
	}
	return err
}

// UpdateMBTIInfo 更新用户MBTI信息
func UpdateMBTIInfo(uid int64, mbtiInfo *model.MBTIInfo) (err error) {
	_, err = dbClient.Exec("update t_user set mbti_str = ? where uid = ?", util.MarshalToStringWithOutErr(mbtiInfo), uid)
	if err != nil {
		log.Errorf("update mbti_info fail,err:%s", err.Error())
		return err
	}
	return err
}

// UpdateUserTagInfo 更新用户标签信息
func UpdateUserTagInfo(uid int64, tagInfo *model.UserTagInfo) (err error) {
	_, err = dbClient.Exec("update t_user set tag_info_str = ? where uid = ?", util.MarshalToStringWithOutErr(tagInfo), uid)
	if err != nil {
		log.Errorf("update tag info fail,err:%s", err.Error())
		return err
	}
	return err
}

// UpdateUserVisible 更新用户的可见性设置
func UpdateUserVisible(uid int64, visible model.Switch) (err error) {
	_, err = dbClient.Exec("update t_user set is_visible = ? where uid = ?", visible, uid)
	if err != nil {
		log.Errorf("update visible fail,err:%s", err.Error())
		return err
	}
	return err
}
