package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

// GetMyGameList 获取我的游戏列表 -- 需要过滤已删除的游戏
func GetMyGameList(uid int64, index string, pageSize int) (list []*model.DbGame, hasMore bool, nextIndex string, err error) {
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	// 获取游戏列表
	list, err = getGameList(" where uid = ? and state != ? order by create_time desc limit ?,?", uid, model.GSDel, offset, pageSize)
	if err != nil {
		log.Errorf("get game list fail,err:%s", err.Error())
		return list, hasMore, nextIndex, err
	}
	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.FormatInt(offset+int64(pageSize), 10)
	}
	return list, hasMore, nextIndex, err
}

// getGameList 获取游戏列表
func getGameList(whereQuery string, params ...interface{}) (list []*model.DbGame, err error) {
	list = make([]*model.DbGame, 0)
	query := "SELECT id, uid, name, prologue, re_total_amount, re_total_num, re_clain_num, answer_list_str, state,version " +
		"UNIX_TIMESTAMP(create_time) AS create_time FROM t_game " + whereQuery
	err = dbClient.FindList(&list, query, params...)
	if err != nil {
		log.Errorf("get game list fail,err: %s", err.Error())
		return list, err
	}
	for i := range list {
		list[i].AnswerList = make([]string, 0)
		_ = jsoniter.UnmarshalFromString(list[i].AnswerListStr, &list[i].AnswerList)
	}
	return list, err
}

// GetGameById 根据id获取游戏信息
func GetGameById(id int64) (game *model.DbGame, err error) {
	list, err := getGameList(" where id = ?", id)
	if err != nil {
		log.Errorf("get game list fail,err: %s", err.Error())
		return game, err
	}
	if len(list) > 0 {
		return list[0], err
	}
	return game, err
}

// AddGame 新增游戏
func AddGame(uid int64, game *model.GameInfo) (effect int64, err error) {
	result, err := dbClient.Exec("INSERT INTO t_game (uid, name, prologue, re_total_amount, re_total_num, answer_list_str) VALUES (?, ?, ?, ?, ?, ?)",
		uid, game.Name, game.Prologue, game.RETotalAmount, game.RETotalNum, util.MarshalToStringWithOutErr(game.AnswerList))
	if err != nil {
		log.Errorf("add game fail,err: %s", err.Error())
		return effect, err
	}
	effect, err = result.RowsAffected()
	if err != nil {
		log.Errorf("get effect fail,err: %s", err.Error())
		return effect, err
	}
	log.Infof("add game success,effect:%d", effect)
	return effect, err
}

// GetGameChatRecordList 获取游戏记录列表
func GetGameChatRecordList(playUid, gameId int64) (list []*model.DbGamePlayRecord, err error) {
	list = make([]*model.DbGamePlayRecord, 0)
	err = dbClient.FindList(&list, "SELECT id, uid, game_id, input, output, result_state, UNIX_TIMESTAMP(create_time) AS create_time "+
		"FROM t_game_play_record where uid = ? and game_id = ? order by create_time asc", playUid, gameId)
	if err != nil {
		log.Errorf("get game play record list fail,err: %s", err.Error())
		return list, err
	}
	return list, err
}

// UpdateGameState 更新游戏状态
func UpdateGameState(gameId, version int64, gameState model.GameState) (effect int64, err error) {
	result, err := dbClient.Exec("update t_game set state = ?, version = version + 1 where id = ? and version = ?", gameState, gameId, version)
	if err != nil {
		log.Errorf("update game state fail,err: %s", err.Error())
		return effect, err
	}
	effect, err = result.RowsAffected()
	if err != nil {
		log.Errorf("get effect fail,err: %s", err.Error())
		return effect, err
	}
	return effect, err
}

// AddGameChatRecord 添加游戏记录失败
func AddGameChatRecord(info *model.DbGamePlayRecord) (err error) {
	_, err = dbClient.Exec("insert into t_game_play_record (uid, game_id, input, output, result_state) VALUES (?, ?, ?, ?, ?)",
		info.Uid, info.GameID, info.Input, info.Output, info.ResultState)
	if err != nil {
		log.Errorf("add game play record fail,err: %s", err.Error())
		return err
	}
	return nil
}

// GameREDecrease 游戏红包减少 -- 并且需要给用户添加红包
func GameREDecrease(uid int64, gameInfo *model.DbGame) (errCode model.ErrCode, errMsg string) {
	// 使用事务进行操作、如果答案正确的话 -- 需要扣减用户红包、并且给玩游戏的用户加上对应的红包金额、添加明细
	err := Transaction(dbClient.DB, func(tx *sqlx.Tx) error {
		// 扣减红包数量，并且如果红包已经发完，则将游戏状态进行变更
		result, errT := tx.Exec("UPDATE t_game SET re_claim_num = re_claim_num + 1, state = CASE WHEN re_claim_num + 1 = re_total_num THEN ? ELSE state END,"+
			"version = CASE WHEN re_claim_num + 1 = re_total_num THEN version + 1 ELSE version END WHERE  id = ? and re_claim_num < re_total_num and version = ?",
			model.GSNoRE, gameInfo.Id, gameInfo.Version)
		if errT != nil {
			log.Errorf("update game record fail,err: %s", errT.Error())
			return errT
		}
		effect, errT := result.RowsAffected()
		if errT != nil {
			log.Errorf("get effect fail,err: %s", errT.Error())
			return errT
		}
		if effect <= 0 {
			errCode = model.ECBan
			errMsg = "游戏信息发生变更，请刷新重试"
			return errT
		}
		// 给用户的资金池添加总计以及详情记录
		// 进行资金池总金额修改
		amount := gameInfo.ReTotalAmount / gameInfo.ReTotalNum // 总金额除以总数量
		errT = fundingPoolAddAmount(tx, uid, model.FPTRe, amount)
		if errT != nil {
			log.Errorf("funding pool add amount err: %s", errT.Error())
			return errT
		}
		// 进行资金池详情信息添加
		errT = fundingPoolAddDetail(tx, &model.DbFundingPoolDetail{
			Uid:      uid,
			PoolType: model.FPTRe,
			Amount:   amount,
		})
		if errT != nil {
			log.Errorf("funding pool add detail err: %s", errT.Error())
			return errT
		}
		return nil
	})
	if err != nil {
		log.Errorf("add user fail,err:%s", err.Error())
		if errCode == model.ErrCodeSuccess {
			errCode = model.ECDbFail
		}
		return errCode, errMsg
	}
	return model.ErrCodeSuccess, ""
}
