package service

import (
	"aDi/dao"
	"aDi/log"
	"aDi/model"
)

// GetUserMapByIdList 获取用户信息map
func GetUserMapByIdList(uidList []int64) (uMap map[int64]*model.UserBaseInfo) {
	uMap = make(map[int64]*model.UserBaseInfo)
	if len(uidList) <= 0 {
		return uMap
	}
	userList, err := dao.GetBUserListById(uidList)
	if err != nil {
		log.Errorf("get user list fail,err:%s", err.Error())
		return
	}
	for i := range userList {
		uMap[userList[i].Uid] = userList[i]
	}
	return uMap
}

// GetUserAllInfo 获取用户信息
func GetUserAllInfo(uid int64) (userAllInfo *model.UserAllInfo, err error) {
	// 获取用户信息
	userInfo, err := dao.GetUserInfoByUid(uid)
	if err != nil {
		log.Errorf("get user info fail,uid:%d,err:%s", uid, err.Error())
		return userAllInfo, err
	}

	// 数据转化
	userAllInfo = model.TransDbUserToAll(userInfo)
	return userAllInfo, err
}

// GetDigitalMapByUid 根据uid列表获取数字人map
func GetDigitalMapByUid(uidList []int64) (dMap map[int64]*model.DbDigitalInfo) {
	dMap = make(map[int64]*model.DbDigitalInfo)
	if len(uidList) <= 0 {
		return dMap
	}
	list, err := dao.GetDigitalListByUid(uidList)
	if err != nil {
		log.Errorf("get digital list fail,err:%s", err.Error())
		return dMap
	}
	for i := range list {
		dMap[list[i].Uid] = list[i]
	}
	return dMap
}
