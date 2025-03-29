package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

// GetMyConversationList 获取我个人和机器人的会话列表 -- 使用会话配置列表作为个人的会话列表
func GetMyConversationList(uid int64, index string, pageSize int) (list []*model.DbConversationConf, hasMore bool, nextIndex string, err error) {
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	list, err = getConversationConfList(" where uid = ? order by create_time desc limit ?,?", uid, offset, pageSize)
	if err != nil {
		log.Errorf("get my conversation list fail,err:%s", err.Error())
		return list, hasMore, nextIndex, err
	}
	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.FormatInt(offset+int64(pageSize), 10)
	}
	return list, hasMore, nextIndex, nil
}

// getConversationConfList 获取会话配置列表 -- 同一个机器人当作一个会话 -- 所以直接使用我与机器人会话配置作为会话列表
func getConversationConfList(whereQuery string, params ...interface{}) (list []*model.DbConversationConf, err error) {
	list = make([]*model.DbConversationConf, 0)
	err = dbClient.FindList(&list, "SELECT id, uid, digital_uid, is_anonymity, chat_conf,last_msg, chat_total_num, "+
		"chat_use_num FROM t_conversation_conf "+whereQuery, params...)
	if err != nil {
		log.Errorf("get my conversation list fail,err:%s", err.Error())
		return list, err
	}
	for i := range list {
		list[i].ChatConfStruct = &model.ConChatConf{}
		if list[i].ChatConf != "" {
			_ = jsoniter.UnmarshalFromString(list[i].ChatConf, list[i].ChatConfStruct)
		}
	}
	return list, err
}

// GetConversationConfById 根据id获取会话配置信息
func GetConversationConfById(id int64) (info *model.DbConversationConf, err error) {
	// 获取列表
	list, err := getConversationConfList(" where id = ? ", id)
	if err != nil {
		log.Errorf("get my conversation conf list fail,err:%s", err.Error())
		return info, err
	}
	if len(list) >= 1 {
		return list[0], err
	}
	return info, err
}

// GetConversationConfByUid 根据uid获取会话配置信息
func GetConversationConfByUid(uid, digitalUid int64) (info *model.DbConversationConf, err error) {
	// 获取列表
	list, err := getConversationConfList(" where uid = ? and digital_uid = ?", uid, digitalUid)
	if err != nil {
		log.Errorf("get conversation conf list fail,err:%s", err.Error())
		return info, err
	}
	if len(list) >= 1 {
		return list[0], err
	}
	return info, err
}

// GetOrAddConversationConfId 获取或者新增对话配置id
func GetOrAddConversationConfId(uid, digitalUid int64) (cID int64) {
	info, err := GetConversationConfByUid(uid, digitalUid)
	if err != nil {
		log.Errorf("get conversation conf id fail,err:%s", err.Error())
		return cID
	}
	if info == nil || info.ID <= 0 {
		// 如果不存在则进行新增
		cID, err = InitConversationConf(uid, digitalUid)
		if err != nil {
			log.Errorf("init conversation conf fail,err:%s", err.Error())
			return cID
		}
		return cID
	}
	return info.ID
}

// InitConversationConf 初始化聊天配置
func InitConversationConf(uid, digitalUid int64) (id int64, err error) {
	result, err := dbClient.Exec("insert into t_conversation_conf (uid, digital_uid, chat_conf)", uid, digitalUid, util.MarshalToStringWithOutErr(&model.BuyChatNumReq{}))
	if err != nil {
		log.Errorf("init conversation conf fail,err:%s", err.Error())
		return id, err
	}
	id, err = result.LastInsertId()
	if err != nil {
		log.Errorf("init conversation conf fail,err:%s", err.Error())
		return id, err
	}
	return id, err
}

// DecreaseConNum 扣减会话聊天次数
func DecreaseConNum(conConfId int64) (effect int64) {
	result, err := dbClient.Exec("update t_conversation_conf set chat_use_num = chat_use_num + 1 where id = ? and chat_total_num > chat_use_num", conConfId)
	if err != nil {
		log.Errorf("decrease con num fail,err:%s", err.Error())
		return effect
	}
	effect, _ = result.RowsAffected()
	return effect
}

// GetMyDigitalConversationList 获取我机器人的会话列表
func GetMyDigitalConversationList(digitalUid int64, index string, pageSize int) (list []*model.DbConversation, hasMore bool, nextIndex string, err error) {
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	list, err = getConversationList(" where digital_uid = ? order by create_time desc limit ?,?", digitalUid, offset, pageSize)
	if err != nil {
		log.Errorf("get my conversation list fail,err:%s", err.Error())
		return list, hasMore, nextIndex, err
	}
	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.FormatInt(offset+int64(pageSize), 10)
	}
	return list, hasMore, nextIndex, err
}

// GetConversationById 根据id获取会话信息
func GetConversationById(id int64) (info *model.DbConversation, err error) { // 获取列表
	list, err := getConversationList(" where id = ? ", id)
	if err != nil {
		log.Errorf("get conversation list fail,err:%s", err.Error())
		return info, err
	}
	if len(list) >= 1 {
		return list[0], err
	}
	return info, err
}

// GetOrAddConversation 获取或者是新增会话--只需要获取id即可
func GetOrAddConversation(uid, digitalUid int64, isAnonymity model.Switch) (id int64) {
	conInfo, err := GetConversationByUid(uid, digitalUid, isAnonymity)
	if err != nil {
		log.Errorf("get conversation info fail,err:%s", err.Error())
		return id
	}
	if conInfo == nil || conInfo.ID <= 0 {
		// 如果会话不存在则进行新增
		id, err = AddConversation(uid, digitalUid, isAnonymity)
		if err != nil {
			log.Errorf("add conversation info fail,err:%s", err.Error())
			return id
		}
		return id
	}
	return conInfo.ID
}

// GetConversationByUid 根据uid获取会话信息
func GetConversationByUid(uid, digitalUid int64, isAnonymity model.Switch) (info *model.DbConversation, err error) {
	// 获取会话信息
	list, err := getConversationList(" where uid = ? and digital_uid = ? and is_anonymity = ?", uid, digitalUid, isAnonymity)
	if err != nil {
		log.Errorf("get conversation list fail,err:%s", err.Error())
		return info, err
	}
	if len(list) >= 1 {
		return list[0], err
	}
	return info, err
}

// AddConversation 添加会话信息
func AddConversation(uid, digitalUid int64, isAnonymity model.Switch) (insertId int64, err error) {
	result, err := dbClient.Exec("INSERT INTO t_conversation (uid, digital_uid, is_anonymity) VALUES (?, ?, ?)", uid, digitalUid, isAnonymity)
	if err != nil {
		log.Errorf("add conversation fail,err:%s", err.Error())
		return insertId, err
	}
	insertId, err = result.LastInsertId()
	if err != nil {
		log.Errorf("get insert id fail,err:%s", err.Error())
		return insertId, err
	}
	return insertId, err
}

// getConversationList 获取会话列表
func getConversationList(whereQuery string, params ...interface{}) (list []*model.DbConversation, err error) {
	list = make([]*model.DbConversation, 0)
	err = dbClient.FindList(&list, "SELECT id, uid, digital_uid, is_anonymity, UNIX_TIMESTAMP(create_time) AS create_time "+
		"FROM t_conversation "+whereQuery, params...)
	if err != nil {
		log.Errorf("get my conversation list fail,err:%s", err.Error())
		return list, err
	}
	return list, err
}

// GetMessageListByUid 根据uid获取消息列表
func GetMessageListByUid(uid, digitalUid int64, index string, pageSize int) (list []*model.DbMessage, hasMore bool, nextIndex string, err error) {
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	list, err = getMessageList(" where uid = ? and digital_uid = ? order by create_time desc limit ?,?", uid, digitalUid, offset, pageSize)
	if err != nil {
		log.Errorf("get my message list fail,err:%s", err.Error())
		return list, hasMore, nextIndex, err
	}
	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.FormatInt(offset+int64(pageSize), 10)
	}
	return list, hasMore, nextIndex, err
}

// GetMessageListByConId 根据会话id获取消息列表
func GetMessageListByConId(conId int64, index string, pageSize int) (list []*model.DbMessage, hasMore bool, nextIndex string, err error) {
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	list, err = getMessageList(" where conversation_id = ? order by create_time desc limit ?,?", conId, offset, pageSize)
	if err != nil {
		log.Errorf("get my message list fail,err:%s", err.Error())
		return list, hasMore, nextIndex, err
	}
	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.FormatInt(offset+int64(pageSize), 10)
	}
	return list, hasMore, nextIndex, err
}

// getMessageList 获取消息列表
func getMessageList(whereQuery string, params ...interface{}) (list []*model.DbMessage, err error) {
	list = make([]*model.DbMessage, 0)
	err = dbClient.FindList(&list, "SELECT id, uid, digital_uid, conversation_id, u_message, d_message, UNIX_TIMESTAMP(create_time) AS create_time "+
		"FROM t_message "+whereQuery, params...)
	if err != nil {
		log.Errorf("get message list fail,err:%s", err.Error())
		return list, err
	}
	return list, err
}

// AddMessage 新增消息
func AddMessage(message *model.DbMessage) (insertId int64, err error) {
	result, err := dbClient.Exec("INSERT INTO t_message (uid, digital_uid, conversation_id, u_message, d_message, chat_conf) VALUES (?, ?, ?, ?, ?, ?)",
		message.UID, message.DigitalUID, message.ConversationID, message.UMessage, message.DMessage, message.ChatConf)
	if err != nil {
		log.Errorf("add message fail,err:%s", err.Error())
		return insertId, err
	}
	insertId, err = result.LastInsertId()
	return insertId, err
}

// UpdateMessageDM 将消息的返回补全
func UpdateMessageDM(id int64, dMessage string) (err error) {
	_, err = dbClient.Exec("update t_message set d_message = ? where id = ?", dMessage, id)
	if err != nil {
		log.Errorf("update message fail,err:%s", err.Error())
		return err
	}
	return err
}

// GetMessageById 根据id获取消息
func GetMessageById(id int64) (message *model.DbMessage, err error) {
	list, err := getMessageList(" where id = ?", id)
	if err != nil {
		log.Errorf("get message list fail,err:%s", err.Error())
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return message, nil
}

// UpsertConversation 更新插入对话配置信息
func UpsertConversation(c *model.DbConversationConf) (effect int64, err error) {
	query := "insert into t_conversation_conf (uid, digital_uid, chat_conf, chat_total_num, is_anonymity) value (?,?,?,?,?)  on duplicate key update chat_conf = values(chat_conf)," +
		"chat_total_num = chat_total_num + values(chat_total_num),is_anonymity = values(is_anonymity)"
	result, err := dbClient.Exec(query, c.UID, c.DigitalUID, c.ChatConf, c.ChatTotalNum, c.IsAnonymity)
	if err != nil {
		log.Errorf("upsert conversation fail,err:%s", err.Error())
		return 0, err
	}
	effect, err = result.RowsAffected()
	if err != nil {
		log.Errorf("upsert conversation fail,err:%s", err.Error())
		return effect, err
	}
	return effect, err
}
