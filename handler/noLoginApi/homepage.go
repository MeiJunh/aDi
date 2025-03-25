package noLoginApi

import (
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"github.com/gin-gonic/gin"
)

// GetOtherHomePageInfoByShareCode 根据分享码获取其他人的主页信息
func (n *NoLoginHandlerImp) GetOtherHomePageInfoByShareCode(c *gin.Context) (rsp model.BaseRsp) {
	shareCode := c.Query("shareCode")
	defer util.TimeCost("GetOtherHomePageByShareCode", shareCode, &rsp.Code)()
	if shareCode == "" {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}

	// 进行分享码拆分
	uid, partOpenid, shareType := service.CommShareCodeSplit(shareCode)
	if uid <= 0 || partOpenid == "" || shareType != service.STHomePage {
		// 解析出来没有uid或open id部分为空字符串，分享类型不对则直接返回报错
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 根据uid获取用户信息
	userInfo, errCode, errMsg := service.ShareCodeCheck(uid, partOpenid)
	if errCode != model.ErrCodeSuccess {
		rsp.WriteMsg(model.CodeMsg{Code: errCode, Msg: errMsg})
		return
	}
	// 获取用户主页配置
	hpConf, err := dao.GetHomepageConfByUid(uid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get homepage conf fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	if hpConf == nil || hpConf.ID <= 0 || hpConf.HpConf == nil {
		// 数据初始化
		hpConf = &model.DbHomepageConf{
			HpConf: &model.HpConf{
				CardList: make([]*model.HpCard, 0),
			},
		}
		return
	}
	// 获取用户的点赞关注统计信息
	statisticInfo, err := dao.GetSocialStatisticInfo(uid)
	if err != nil {
		// 没有获取到也不用报错，降级处理
		log.Errorf("get social statistic info fail, uid:%d, err:%s", uid, err.Error())
	}
	if statisticInfo == nil {
		statisticInfo = &model.DBFollowStatisticInfo{}
	}
	digitalInfo, err := dao.GetDigitalInfo(uid)
	if err != nil {
		log.Errorf("get digital info fail, uid:%d, err:%s", uid, err.Error())
	}
	if digitalInfo == nil {
		digitalInfo = &model.DbDigitalInfo{
			DigitalName: userInfo.Nick, // 使用用户昵称进行降级
		}
	}
	rsp.Data = &model.GetOtherHomePageInfoByShareCodeRsp{
		HpConfInfo: hpConf.HpConf,
		StatisticInfo: &model.FollowStatisticInfo{
			FollowNum: statisticInfo.FollowNum,
			FavorNum:  statisticInfo.FavorNum,
			ViewNum:   statisticInfo.ViewNum,
		},
		DigitalInfo: &model.DigitalShadowInfo{
			UserInfo: &model.UserAllInfo{
				UserBaseInfo: model.UserBaseInfo{
					Nick: userInfo.Nick,
					Icon: userInfo.Icon,
					Age:  userInfo.Age,
					Sex:  userInfo.Sex,
				},
				TagInfo: model.TransDbUserToAll(userInfo).TagInfo,
			},
			DigitalName: digitalInfo.DigitalName,
		},
	}
	return
}
