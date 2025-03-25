package service

import (
	"aDi/dao"
	"aDi/log"
	"aDi/model"
)

// GetDigitalShadowInfoByUid 根据uid获取数字人信息
func GetDigitalShadowInfoByUid(uid int64) (digInfo *model.DigitalShadowInfo, errCode model.ErrCode, errMsg string) {
	// 获取个人信息
	userInfo, err := GetUserAllInfo(uid)
	if err != nil {
		log.Errorf("get user info fail,err：%s", err.Error())
		return digInfo, model.ErrIDbFail.Code, model.ErrIDbFail.Msg
	}
	if userInfo == nil || userInfo.Uid <= 0 {
		return digInfo, model.ECNoExist, "用户不存在"
	}
	// 获取自己的数字人信息
	digitalInfo, err := dao.GetDigitalInfo(uid)
	if err != nil {
		log.Errorf("get digital info fail,err:%s", err.Error())
		return digInfo, model.ErrIDbFail.Code, model.ErrIDbFail.Msg
	}
	// 数据转化
	digInfo = model.DbDigitalInfoTrans(digitalInfo)
	// 个人信息填充
	digInfo.UserInfo = userInfo
	return digInfo, model.ErrCodeSuccess, ""
}

// GetDigitalInfoByUid 根据uid获取机器人信息
func GetDigitalInfoByUid(uid int64) (info *model.DigitalShadowInfo, errCode model.ErrCode, errMsg string) {
	digitalInfo, err := dao.GetDigitalInfo(uid)
	if err != nil {
		log.Errorf("get digital info fail,err:%s", err.Error())
		return info, model.ErrIDbFail.Code, model.ErrIDbFail.Msg
	}
	// 数据转化
	info = model.DbDigitalInfoTrans(digitalInfo)
	return info, model.ErrCodeSuccess, ""
}
