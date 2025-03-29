package service

import (
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"time"
)

const (
	GameMaxLen = 5 // 游戏最大轮数为5轮
)

// GetGameStateByRecordList 根据聊天列表以及当前游戏状态获取指定用户的游戏状态
func GetGameStateByRecordList(game *model.DbGame, recordList []*model.GameChatInfo) (rState model.GameState) {
	if game.State != model.GSDefault {
		// 状态不正常直接返回
		return game.State
	}
	// 游戏过期单独判断
	if time.Now().Unix()-game.CreateTime >= util.DaySecond {
		return model.GSExpire
	}
	for i := range recordList {
		if recordList[i].ResultState == model.GCSRight {
			return model.GSPlayWin
		}
	}
	if len(recordList) >= GameMaxLen {
		return model.GSPlayOver
	}
	return game.State
}

// GetGameRecordList 获取游戏记录列表
func GetGameRecordList(playUid, gameId int64) (list []*model.GameChatInfo) {
	list = make([]*model.GameChatInfo, 0)
	dList, err := dao.GetGameChatRecordList(playUid, gameId)
	if err != nil {
		log.Errorf("get game records fail, err: %s", err.Error())
		return list
	}
	for i := range dList {
		list = append(list, &model.GameChatInfo{
			Input:       dList[i].Input,
			Output:      dList[i].Output,
			ResultState: dList[i].ResultState,
			CreateTime:  dList[i].CreateTime,
		})
	}
	return list
}

// DelGame 删除游戏
func DelGame(gameInfo *model.DbGame) (errCode model.ErrCode, errMsg string) {
	// 删除游戏
	effect, err := dao.UpdateGameState(gameInfo.Id, gameInfo.Version, model.GSDel)
	if err != nil {
		log.Errorf("update game state fail, err: %s", err.Error())
		return model.ECDbFail, "删除游戏失败，请稍后再试"
	}
	if effect <= 0 {
		// 游戏状态发生变化 -- 返回不让操作
		return model.ECBan, "游戏发生变化，请刷新重试"
	}
	if gameInfo.State == model.GSDefault && gameInfo.ReClaimNum < gameInfo.ReTotalNum {
		// 如果游戏是正常状态--则需要将剩余的红包金额返回 -- todo
	}
	return model.ErrCodeSuccess, ""
}

// GetGameInfoByShareCode 根据分享码获取游戏信息
func GetGameInfoByShareCode(uid int64, shareCode string) (gameInfo *model.DbGame, gamePlayInfo *model.GamePlayInfo, errCode model.ErrCode, errMsg string) {
	var err error
	// 根据分享码拆分出游戏id
	dUid, pOpenId, shareType, gameId := IdShareCodeSplit(shareCode)
	if shareType != STGame {

		return gameInfo, gamePlayInfo, model.ECBan, "该邀请码不是游戏邀请码"
	}
	// 防伪码校验
	_, errCode, errMsg = ShareCodeCheck(dUid, pOpenId)
	if errCode != model.ErrCodeSuccess {
		return gameInfo, gamePlayInfo, errCode, errMsg
	}
	// 获取游戏信息
	gameInfo, err = dao.GetGameById(gameId)
	if err != nil {
		log.Errorf("get game info fail,err:%s", err.Error())
		return gameInfo, gamePlayInfo, model.ErrIDbFail.Code, model.ErrIDbFail.Msg
	}
	if gameInfo == nil || gameInfo.Id <= 0 {
		return gameInfo, gamePlayInfo, model.ECNoExist, "当前游戏不存在"
	}

	// 查看该用户与这个游戏的对话记录
	chatList := GetGameRecordList(uid, gameId)
	// 获取数字人信息 -- 返回icon和name
	digitalInfo, _ := dao.GetDigitalInfo(gameInfo.Uid)
	gamePlayInfo = &model.GamePlayInfo{
		GameName:  gameInfo.Name,
		Prologue:  gameInfo.Prologue,
		ChatList:  chatList,
		GameState: GetGameStateByRecordList(gameInfo, chatList),
		DigitalInfo: &model.ChatDigitalInfo{
			Name: digitalInfo.GetDigitalName(),
			Icon: digitalInfo.GetIcon(),
		},
	}
	return gameInfo, gamePlayInfo, errCode, errMsg
}

// ChatWithGame 玩游戏
func ChatWithGame(uid int64, input string, gameInfo *model.DbGame) (errCode model.ErrCode, errMsg string) {
	var err error
	// 获取聊天id -- 使用不匿名的聊天
	cId := dao.GetOrAddConversation(uid, gameInfo.Uid, model.SwitchOff)
	if cId <= 0 {
		log.Error("get conversation fail")
		return model.ECDbFail, "获取聊天会话失败"
	}
	// 进行游戏结果判断
	gameResult := model.GCSWrong
	output := "抱歉答错了"
	if util.InSliceStr(gameInfo.AnswerList, input) {
		gameResult = model.GCSRight
		output = "恭喜答对了"
		// 使用事务进行操作、如果答案正确的话 -- 需要扣减用户红包、并且给玩游戏的用户加上对应的红包金额、添加明细
		errCode, errMsg = dao.GameREDecrease(uid, gameInfo)
		if errCode != model.ErrCodeSuccess {
			log.Errorf("game re decrease fail, err code:%d,err msg:%s", errCode, errMsg)
			return errCode, errMsg
		}
	}
	// 添加游戏记录
	err = dao.AddGameChatRecord(&model.DbGamePlayRecord{
		Uid:         uid,
		GameID:      gameInfo.Id,
		Input:       input,
		Output:      output,
		ResultState: gameResult,
	})
	if err != nil {
		log.Errorf("add game records fail, err: %s", err.Error())
		return model.ECDbFail, "添加游戏记录失败"
	}
	// 聊天记录也增加
	_, err = dao.AddMessage(&model.DbMessage{
		UID:            uid,
		DigitalUID:     gameInfo.Uid,
		ConversationID: cId,
		UMessage:       input,
		DMessage:       output,
	})
	if err != nil {
		log.Errorf("add message fail, err: %s", err.Error())
		return model.ECDbFail, "添加对话记录失败"
	}
	return model.ErrCodeSuccess, ""
}
