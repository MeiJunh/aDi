package loginApi

import (
	"aDi/dao"
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetChatRankList 获取用户的排行榜单列表
func (l *LoginHandlerImp) GetChatRankList(c *gin.Context) (rsp model.BaseRsp) {
	// 获取查看人的榜单列表 -- 目前只有聊天次数的榜单
	dUidStr := c.Query("digitalUid")
	digitalUid := util.ToInt64(dUidStr)
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetChatRankList", dUidStr, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 || digitalUid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// TODO 验证用户和被查看人的关系
	rankList, err := dao.GetChatTopList(digitalUid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get rank list fail,err:%s", err.Error())
		return
	}
	if len(rankList) == 0 {
		// 设置默认值然后直接返回
		rsp.Data = &model.GetListRsp{
			List: []*model.RankInfo{},
		}
		return
	}
	// 获取用户信息进行补全
	uidList := make([]int64, 0, len(rankList))
	for i := range rankList {
		uidList = append(uidList, rankList[i].ChatUid)
	}
	// 获取用户信息map
	uMap := service.GetUserMapByIdList(uidList)
	rList := make([]*model.RankInfo, 0)
	for i := range rankList {
		uBase := uMap[rankList[i].ChatUid]
		if rankList[i].IsAnonymity == model.SwitchOn || uBase == nil {
			uBase = &model.UserBaseInfo{
				Nick: "张三",
				// TODO 填充默认头像
			}
		}

		rList = append(rList, &model.RankInfo{
			UserBaseInfo: *uBase,
			Score:        rankList[i].ChatNum,
		})
	}
	rsp.Data = &model.GetListRsp{
		List: rList,
	}
	return
}
