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
用户信息相关
*/

// GetUserInfo 获取用户个人资料页信息
func (l *LoginHandlerImp) GetUserInfo(c *gin.Context) (rsp model.BaseRsp) {
	// 获取uid
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetUserInfo", uid, &rsp.Code)()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return rsp
	}

	// 获取用户信息
	userInfo, err := service.GetUserAllInfo(uid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Error("get user info fail,uid:%d,err:%s", uid, err.Error())
		return rsp
	}
	rsp.Data = userInfo
	return
}

// ModUserInfo 修改用户个人资料页信息
func (l *LoginHandlerImp) ModUserInfo(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.UserAllInfo{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ModUserInfo", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return rsp
	}
	if req.Uid > 0 && req.Uid != uid {
		rsp.WriteMsg(model.CodeMsg{
			Code: model.ECBan,
			Msg:  "不能修改其他人的信息",
		})
		return
	}
	// todo 进行填写参数校验
	// 进行信息修改
	err = dao.UpdateUserBaseInfo(req)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Error("update user info fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	return
}

// ModUserMBTI 修改用户的MBTI信息
func (l *LoginHandlerImp) ModUserMBTI(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.MBTIInfo{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ModUserMBTI", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return rsp
	}

	// 更新MBTI信息
	err = dao.UpdateMBTIInfo(uid, req)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Error("update user mbti_info fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	return
}

// ModUserTagInfo 修改用户标签信息
func (l *LoginHandlerImp) ModUserTagInfo(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.UserTagInfo{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ModUserTagInfo", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return rsp
	}

	// 更新标签信息
	err = dao.UpdateUserTagInfo(uid, req)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Error("update user mbti_info fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	return
}

// ModVisible 设置别人是否可见
func (l *LoginHandlerImp) ModVisible(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.ModVisibleReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("ModVisible", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return rsp
	}

	// 更新用户的可见性设置
	err = dao.UpdateUserVisible(uid, req.Visible)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Error("update user visible fail,uid:%d,err:%s", uid, err.Error())
		return
	}
	return
}
