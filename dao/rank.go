package dao

import (
	"aDi/log"
	"aDi/model"
)

// 聊天统计

// GetChatTopList 获取用户机器人聊天排行榜 -- 目前只用获取前10
func GetChatTopList(digitalUid int64) (list []*model.DbChatStatistic, err error) {
	list = make([]*model.DbChatStatistic, 0)
	err = dbClient.FindList(&list, "select digital_uid,chat_uid,is_anonymity,chat_num from t_chat_statistic where digital_uid = ?"+
		" order by chat_num desc limit 10", digitalUid)
	if err != nil {
		log.Errorf("get chat top list fail,err: %s", err.Error())
		return list, err
	}
	return list, err
}
