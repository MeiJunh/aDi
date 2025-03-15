package loginApi

import (
	"aDi/dao"
	"aDi/handler/comm"
	"aDi/log"
	"aDi/model"
	"aDi/service"
	"aDi/util"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

/*
对话、聊天
*/

// GetMyConversationList 获取我的会话列表 -- 用户获取自己和机器人的会话聊天列表
func (l *LoginHandlerImp) GetMyConversationList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyConversationList", uid, &rsp.Code, strconv.FormatInt(uid, 10))()
	if uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 获取个人和机器人的会话列表
	cList, hasMore, nextIndex, err := dao.GetMyConversationList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get my conversation list fail,err:%s", err.Error())
		return
	}
	// 数据信息转化
	list := make([]*model.Conversation, 0)
	uidList := make([]int64, 0)
	for i := range cList {
		uidList = append(uidList, cList[i].DigitalUID)
		list = append(list, &model.Conversation{
			Id:      cList[i].ID,
			LastMsg: cList[i].LastMsg,
			DigitalChatInfo: model.ChatDigitalInfo{
				Uid:         cList[i].DigitalUID,
				ConChatConf: *cList[i].ChatConfStruct,
			},
		})
	}
	// 获取数字人信息map
	dMap := service.GetDigitalMapByUid(uidList)
	for i := range list {
		// 进行数字人信息填充
		list[i].DigitalChatInfo.Name = dMap[list[i].DigitalChatInfo.Uid].GetDigitalName()
		list[i].DigitalChatInfo.Icon = dMap[list[i].DigitalChatInfo.Uid].GetIcon()
	}
	rsp.Data = &model.GetListRsp{
		List:      []*model.Conversation{},
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// GetMyMessageList 获取自己某次会话的聊天记录 -- 从最新的往前面获取
func (l *LoginHandlerImp) GetMyMessageList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetMessageListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyMessageList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.ConversationId <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 先根据传过来的会话id获取会话配置信息
	cConf, err := dao.GetConversationConfById(req.ConversationId)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get my conversation list fail,err:%s", err.Error())
		return
	}
	if cConf == nil || cConf.ID <= 0 || cConf.UID != uid {
		rsp.WriteMsg(model.CodeMsg{Code: model.ECNoExist, Msg: "对应的会话不存在"})
		log.Errorf("conversation not exist,cConf:%+v", cConf)
		return
	}
	// 然后根据当前uid和机器人uid获取他们之间所有的消息列表
	list, hasMore, nextIndex, err := dao.GetMessageListByUid(uid, cConf.DigitalUID, req.Index, req.PageSize)
	if err != nil {
		log.Errorf("get my message list fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIDbFail)
		return
	}
	rList := make([]*model.ChatMessage, 0)
	for i := range list {
		rList = append(rList, &model.ChatMessage{
			Id:         list[i].ID,
			UMessage:   list[i].UMessage,
			DMessage:   list[i].DMessage,
			CreateTime: list[i].CreateTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      []*model.ChatMessage{},
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// GetMyDigitalConList 获取自己数字分身的会话记录列表
func (l *LoginHandlerImp) GetMyDigitalConList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyDigitalConList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	// 获取数字人会话列表
	cList, hasMore, nextIndex, err := dao.GetMyDigitalConversationList(uid, req.Index, req.PageSize)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get my digital conversation list fail,err:%s", err.Error())
		return
	}
	uidList := make([]int64, 0)
	rList := make([]*model.DigitalConversation, 0)
	for i := range cList {
		if cList[i].IsAnonymity != model.SwitchOn {
			uidList = append(uidList, cList[i].UID)
		}
	}
	uMap := service.GetUserMapByIdList(uidList)
	for i := range cList {
		uBase := uMap[cList[i].UID]
		if cList[i].IsAnonymity == model.SwitchOn || uBase == nil {
			uBase = &model.UserBaseInfo{
				Nick: "张三",
				// TODO 填充默认头像
			}
		}
		rList = append(rList, &model.DigitalConversation{
			Id:               cList[i].ID,
			CharUserBaseInfo: uBase,
			IsAnonymity:      cList[i].IsAnonymity,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// GetMyDigitalMessageList 获取自己数字机器人某次会话的聊天记录 -- 从最新的往前面获取
func (l *LoginHandlerImp) GetMyDigitalMessageList(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.GetMessageListReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("GetMyDigitalMessageList", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.ConversationId <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 获取会话信息是否存在以及匹配
	cInfo, err := dao.GetConversationById(req.ConversationId)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get conversation info fail,err:%s", err.Error())
		return
	}
	if cInfo == nil || cInfo.ID <= 0 || cInfo.DigitalUID != uid {
		rsp.WriteMsg(model.CodeMsg{Code: model.ECNoExist, Msg: "对应的会话不存在"})
		log.Errorf("conversation not exist,cInfo:%+v", cInfo)
		return
	}
	// 获取消息列表
	list, hasMore, nextIndex, err := dao.GetMessageListByConId(req.ConversationId, req.Index, req.PageSize)
	if err != nil {
		log.Errorf("get my message list fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIDbFail)
		return
	}
	rList := make([]*model.ChatMessage, 0)
	for i := range list {
		rList = append(rList, &model.ChatMessage{
			Id:         list[i].ID,
			UMessage:   list[i].UMessage,
			DMessage:   list[i].DMessage,
			CreateTime: list[i].CreateTime,
		})
	}
	rsp.Data = &model.GetListRsp{
		List:      rList,
		HasMore:   hasMore,
		NextIndex: nextIndex,
	}
	return
}

// AiChatGenerate ai对话问答生成
func (l *LoginHandlerImp) AiChatGenerate(c *gin.Context) (rsp model.BaseRsp) {
	req := &model.AiChatGenerateReq{}
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("AiChatGenerate", req, &rsp.Code, strconv.FormatInt(uid, 10))()
	_, err := comm.ReadBodyFromGin(c, req)
	if err != nil {
		log.Errorf("read from body fail,err:%s", err.Error())
		rsp.WriteMsg(model.ErrIParse)
		return
	}
	if uid <= 0 || req.DigitalUid <= 0 || req.Content == "" {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}
	// 先根据数字人id获取对应的会话配置
	conConf, err := dao.GetConversationConfByUid(uid, req.DigitalUid)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get conversation info fail,err:%s", err.Error())
		return
	}
	if conConf == nil || conConf.ID <= 0 || conConf.ChatUseNum >= conConf.ChatTotalNum {
		// 没有对应的会话配置 -- 不让进行聊天生成
		rsp.WriteMsg(model.CodeMsg{Code: model.ECBan, Msg: "没有对话次数"})
		return
	}
	// 获取或者新增会话信息
	cId := dao.GetOrAddConversation(uid, req.DigitalUid, conConf.IsAnonymity)
	if cId <= 0 {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("conversation not exist,cId:%d", cId)
		return
	}
	// 异步进行对话生成
	asyncResult := &model.AiAsyncResult{}
	// 添加异步任务
	keyId := dao.AddAsyncInfo(uid, util.MarshalToStringWithOutErr(asyncResult))
	if keyId <= 0 {
		rsp.WriteMsg(model.CodeMsg{Code: model.ECDbFail, Msg: "添加任务失败,请刷新重试"})
		return
	}
	// 添加消息
	messageId, err := dao.AddMessage(&model.DbMessage{
		UID:            uid,
		DigitalUID:     req.DigitalUid,
		ConversationID: cId,
		UMessage:       req.Content,
		ChatConf:       conConf.ChatConf,
	})
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("add message fail,err:%s", err.Error())
		return
	}
	// 进行次数扣减
	effect := dao.DecreaseConNum(conConf.ID)
	if effect <= 0 {
		// 会话次数不够
		rsp.WriteMsg(model.CodeMsg{Code: model.ECBan, Msg: "没有对话次数"})
		return
	}
	util.GoSafe(func() {
		defer func() {
			// 最后结果写入 -- 包括最后全部成功或者是失败写入
			_ = dao.UpdateAsyncInfo(keyId, util.MarshalToStringWithOutErr(asyncResult))
			// 如果结果错误则返回次数
			if asyncResult.AiStatus != model.ErrCodeSuccess2 {
			}
		}()
		// 根据对应的数字人设定进行参数获取
		aiReq := service.GetChatAiReq(req.Content, nil)
		// 调用ai进行对话生成
		message, _, errCode, errMsg := service.AiJsonContentGenerate(aiReq, keyId, asyncResult)
		if errCode != model.ErrCodeSuccess {
			asyncResult.AiStatus = errCode
			asyncResult.AiErrMsg = string(errMsg)
			log.Errorf("chat create fail,err code:%d,err msg:%s", errCode, errMsg)
			return
		}
		asyncResult.AiStatus = model.ErrCodeSuccess2
		asyncResult.AiResult = message
		// 将对话中的回答补全
		_ = dao.UpdateMessageDM(messageId, message)
		return
	})
	// 调用接口进行数字人问题回答
	rsp.Data = &model.AiChatGenerateRsp{
		KeyCode:   strconv.FormatInt(keyId, 10),
		MessageId: messageId,
	}
	return
}

// AiChatResultGet ai对话结果异步获取
func (l *LoginHandlerImp) AiChatResultGet(c *gin.Context) (rsp model.BaseRsp) {
	// 通过param 获取 keyCode
	keyCode := c.Query("keyCode")
	keyId := util.ToInt64(keyCode)
	uid := comm.GetUidFromCon(c)
	defer util.TimeCost("AiChatResultGet", keyCode, &rsp.Code, strconv.FormatInt(uid, 10))()
	if keyId <= 0 || uid <= 0 {
		rsp.WriteMsg(model.ErrIInvalidParam)
		return
	}

	// 获取异步结果信息
	asyncInfo, err := dao.GetAsyncInfoById(keyId)
	if err != nil {
		rsp.WriteMsg(model.ErrIDbFail)
		log.Errorf("get async info fail,err:%s", err.Error())
		return
	}
	if asyncInfo == nil || asyncInfo.ID <= 0 {
		rsp.Data = &model.AiAsyncResult{AiStatus: model.ErrCodeSuccess} // 需要继续轮询
		return
	}
	// 进行结果解析
	r := &model.AiAsyncResult{}
	_ = jsoniter.UnmarshalFromString(asyncInfo.Result, r)
	rsp.Data = r
	return
}
