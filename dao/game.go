package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"strconv"
)

// GetMyGameList 获取我的游戏列表
func GetMyGameList(uid int64, index string, pageSize int) (list []*model.DbGame, hasMore bool, nextIndex string, err error) {
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	// 获取游戏列表
	list, err = getGameList(" where uid = ? order by create_time desc limit ?,?", uid, offset, pageSize)
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
	query := "SELECT id, uid, name, prologue, re_total_amount, re_total_num, re_clain_num, answer_list_str, state, " +
		"UNIX_TIMESTAMP(create_time) AS create_time FROM t_game " + whereQuery
	err = dbClient.FindList(&list, query, params...)
	if err != nil {
		log.Errorf("get game list fail,err: %s", err.Error())
		return list, err
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
