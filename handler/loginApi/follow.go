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

/*
点赞、关注相关
*/

// GetMyFollowStatisticInfo 获取我自己点赞关注的统计信息
func (l *LoginHandlerImp) GetMyFollowStatisticInfo(c *gin.Context) (rsp model.BaseRsp) {
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyFollowList", uid, &rsp.Code)()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取自己的点赞关注统计信息
	statisticInfo, err := dao.GetSocialStatisticInfo(uid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get statistic info fail,err:%s", err.Error())
		return
	}
	if statisticInfo == nil || statisticInfo.Id <= 0 {
		statisticInfo = &model.DBFollowStatisticInfo{}
		// 如果没有统计信息，则进行统计信息初始化
		dao.InitSocialStatisticInfo(uid)
	}
	rsp.Data = &model.FollowStatisticInfo{
		FollowNum: statisticInfo.FollowNum,
		FavorNum:  statisticInfo.FavorNum,
		ViewNum:   statisticInfo.ViewNum,
	}
	return
}

// GetMyFollowList 获取我的关注列表
func (l *LoginHandlerImp) GetMyFollowList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyFollowList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取我的关注列表，follower是我
	list, hasMore, nextIndex, err := dao.GetMyFollowList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get my follow list fail,err:%s", err.Error())
		return
	}
	// 填充用户信息
	uidList := make([]int64, 0)
	for i := range list {
		uidList = append(uidList, list[i].Uid)
	}
	// 获取用户信息map
	uMap := service.GetUserMapByIdList(uidList)
	rList := make([]*model.FollowInfo, 0)
	for i := range list {
		userTmp := uMap[list[i].Uid]
		if userTmp == nil {
			userTmp = &model.UserBaseInfo{Uid: list[i].Uid}
		}
		rList = append(rList, &model.FollowInfo{
			UserBaseInfo: *userTmp,
			FollowTime:   list[i].CTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// AddFollow 添加关注
func (l *LoginHandlerImp) AddFollow(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.FollowReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("AddFollow", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.Uid <= 0 {
		// 参数不满足直接返回
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// todo 判断别人是否是可见的或者自己是否有通过分享获取过信息
	// 添加自己对别人的关注
	dao.AddFollowInfo(req.Uid, uid)
	return
}

// CancelFollow 取消关注
func (l *LoginHandlerImp) CancelFollow(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.FollowReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("CancelFollow", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.Uid <= 0 {
		// 参数不满足直接返回
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}

	// 取消关注 -- 取消自己对别人的关注
	dao.CancelFollow(req.Uid, uid)
	return
}

// GetFollowMeList 获取关注我的列表
func (l *LoginHandlerImp) GetFollowMeList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyFollowList", req, &rsp.Code)()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 获取关注我的列表，uid是我
	list, hasMore, nextIndex, err := dao.GetFollowMeList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get follow me list fail,err:%s", err.Error())
		return
	}
	// 填充用户信息
	uidList := make([]int64, 0)
	for i := range list {
		uidList = append(uidList, list[i].Follower)
	}
	// 获取用户信息map
	uMap := service.GetUserMapByIdList(uidList)
	rList := make([]*model.FollowInfo, 0)
	for i := range list {
		userTmp := uMap[list[i].Follower]
		if userTmp == nil {
			userTmp = &model.UserBaseInfo{Uid: list[i].Follower}
		}
		rList = append(rList, &model.FollowInfo{
			UserBaseInfo: *userTmp,
			FollowTime:   list[i].CTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// AddFavor 添加点赞
func (l *LoginHandlerImp) AddFavor(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.FollowReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("AddFavor", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.Uid <= 0 {
		// 参数不满足直接返回
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// todo 判断别人是否是可见的或者自己是否有通过分享获取过信息
	// 添加对别人的点赞
	dao.AddFavorInfo(req.Uid, uid)
	return
}

// GetFavorMeList 获取点赞我的列表
func (l *LoginHandlerImp) GetFavorMeList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetFavorMeList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 获取点赞我的列表，uid是我
	list, hasMore, nextIndex, err := dao.GetFavorMeList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get favor me list fail,err:%s", err.Error())
		return
	}
	// 填充用户信息
	uidList := make([]int64, 0)
	for i := range list {
		uidList = append(uidList, list[i].Liker)
	}
	// 获取用户信息map
	uMap := service.GetUserMapByIdList(uidList)
	rList := make([]*model.FollowInfo, 0)
	for i := range list {
		userTmp := uMap[list[i].Liker]
		if userTmp == nil {
			userTmp = &model.UserBaseInfo{Uid: list[i].Liker}
		}
		rList = append(rList, &model.FollowInfo{
			UserBaseInfo: *userTmp,
			FollowTime:   list[i].CTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// GetMyViewList 获取我浏览过的列表
func (l *LoginHandlerImp) GetMyViewList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyViewList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 获取我浏览的列表，viewer是我
	list, hasMore, nextIndex, err := dao.GetMyViewList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get view me list fail,err:%s", err.Error())
		return
	}
	// 填充用户信息
	uidList := make([]int64, 0)
	for i := range list {
		uidList = append(uidList, list[i].Viewer)
	}
	// 获取用户信息map
	uMap := service.GetUserMapByIdList(uidList)
	rList := make([]*model.FollowInfo, 0)
	for i := range list {
		userTmp := uMap[list[i].Viewer]
		if userTmp == nil {
			userTmp = &model.UserBaseInfo{Uid: list[i].Viewer}
		}
		rList = append(rList, &model.FollowInfo{
			UserBaseInfo: *userTmp,
			FollowTime:   list[i].CTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// GetViewMeList 获取浏览过我的列表
func (l *LoginHandlerImp) GetViewMeList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetViewMeList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 获取浏览我的列表，uid是我
	list, hasMore, nextIndex, err := dao.GetViewMeList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get view me list fail,err:%s", err.Error())
		return
	}
	// 填充用户信息
	uidList := make([]int64, 0)
	for i := range list {
		uidList = append(uidList, list[i].Viewer)
	}
	// 获取用户信息map
	uMap := service.GetUserMapByIdList(uidList)
	rList := make([]*model.FollowInfo, 0)
	for i := range list {
		userTmp := uMap[list[i].Viewer]
		if userTmp == nil {
			userTmp = &model.UserBaseInfo{Uid: list[i].Viewer}
		}
		rList = append(rList, &model.FollowInfo{
			UserBaseInfo: *userTmp,
			FollowTime:   list[i].CTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}
