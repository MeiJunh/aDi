package service

import (
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"fmt"
	"math/rand"
	"strings"
)

// ShareType 分享类型
type ShareType int32

const (
	STHomePage ShareType = 1 // 分享主页
	STGame     ShareType = 2 // 分享游戏
	STDigital  ShareType = 3 // 分享数字人

	CommShareCodeFormat   = "%d_%d_%s"                  // uid_shareType_partOfOpenId
	ShareCodeWithIdFormat = CommShareCodeFormat + "_%d" // uid_shareType_partOfOpenId_id
)

// GetShareCodeByUidAndId 根据uid和分享类型、id获取分享码
func GetShareCodeByUidAndId(uid int64, openId string, shareType ShareType, id int64) (shareCode string) {
	// 分享的话根据uid+open id的一部分作为相互校验，然后加上分享类型
	shareCode = fmt.Sprintf(ShareCodeWithIdFormat, uid, shareType, RandomSubstring(openId), id)
	return shareCode
}

// GetShareCodeByUid 根据uid和分享类型获取分享码
func GetShareCodeByUid(uid int64, openId string, shareType ShareType) (shareCode string) {
	// 分享的话根据uid+open id的一部分作为相互校验，然后加上分享类型
	shareCode = fmt.Sprintf(CommShareCodeFormat, uid, shareType, RandomSubstring(openId))
	return shareCode
}

// RandomSubstring 随机截取字符串，长度在5到10之间
func RandomSubstring(input string) string {
	// 检查输入字符串长度是否足够
	if len(input) < 5 {
		// 直接返回
		return input
	}

	// 定义截取长度的范围
	minLength := 5
	maxLength := 10

	// 确定最大允许的截取长度
	// 如果字符串长度小于最大长度，则最大截取长度为字符串长度
	maxAllowedLength := maxLength
	if len(input) < maxLength {
		maxAllowedLength = len(input)
	}

	// 随机生成截取的长度
	length := rand.Intn(maxAllowedLength-minLength) + minLength

	// 计算起始位置的最大值
	// 起始位置 + 截取长度不能超过字符串长度
	// 随机生成起始位置
	start := rand.Intn(len(input) - length)

	// 截取并返回子字符串
	return input[start : start+length]
}

// CommShareCodeSplit 普通分享码拆分
func CommShareCodeSplit(shareCode string) (uid int64, partOpenId string, shareType ShareType) {
	// 普通分享码由uid share type与part of open id构成
	sli := strings.Split(shareCode, "_")
	if len(sli) <= 2 {
		// 拆分出来的数据少于等于2个直接返回
		return uid, partOpenId, shareType
	}
	uid = util.ToInt64(sli[0])
	shareType = ShareType(util.ToInt64(sli[1]))
	partOpenId = sli[2]
	return uid, partOpenId, shareType
}

// IdShareCodeSplit 带id的分享码拆分
func IdShareCodeSplit(shareCode string) (uid int64, partOpenId string, shareType ShareType, id int64) {
	// 普通分享码由uid share type与part of open id,id构成
	sli := strings.Split(shareCode, "_")
	if len(sli) <= 3 {
		// 拆分出来的数据少于等于3个直接返回
		return uid, partOpenId, shareType, id
	}
	uid = util.ToInt64(sli[0])
	shareType = ShareType(util.ToInt64(sli[1]))
	partOpenId = sli[2]
	id = util.ToInt64(sli[3])
	return uid, partOpenId, shareType, id
}

// ShareCodeCheck 邀请码校验
func ShareCodeCheck(uid int64, partOpenid string) (userInfo *model.DBUserInfo, errCode model.ErrCode, errMsg string) {
	// 根据uid获取用户信息
	var err error
	userInfo, err = dao.GetUserInfoByUid(uid)
	if err != nil {
		log.Errorf("get user info by uid fail, uid: %d, err: %s", uid, err.Error())
		return userInfo, model.ErrIDbFail.Code, model.ErrIDbFail.Msg
	}
	if userInfo == nil || userInfo.Uid <= 0 {
		return userInfo, model.ECNoExist, "对应的用户不存在"
	}
	if !strings.Contains(userInfo.OpenID, partOpenid) {
		// uid和open id对应不上

		return userInfo, model.ErrIInvalidParam.Code, model.ErrIInvalidParam.Msg
	}
	return userInfo, model.ErrCodeSuccess, ""
}
