package dao

import (
	"aDi/log"
	"aDi/model"
)

// AddAsyncInfo 添加异步信息
func AddAsyncInfo(uid int64, result string) (insertId int64) {
	r, err := dbClient.Exec("INSERT INTO t_async_result (uid, result) VALUES (?, ?)", uid, result)
	if err != nil {
		log.Errorf("add async info fail,err: %s", err.Error())
		return insertId
	}
	insertId, err = r.LastInsertId()
	if err != nil {
		log.Errorf("add async info fail,err: %s", err.Error())
		return insertId
	}
	return insertId
}

// UpdateAsyncInfo 更新异步信息
func UpdateAsyncInfo(id int64, result string) (err error) {
	_, err = dbClient.Exec("UPDATE t_async_result SET result = ? WHERE id = ?", result, id)
	if err != nil {
		log.Errorf("update async info fail,err: %s", err.Error())
		return err
	}
	return err
}

// GetAsyncInfoById 获取异步信息
func GetAsyncInfoById(id int64) (result *model.DbAsyncResult, err error) {
	result = &model.DbAsyncResult{}
	err = dbClient.FindOneWithNull(result, "SELECT id, uid, result FROM t_async_result where id = ?", id)
	if err != nil {
		log.Errorf("get async info fail,err: %s", err.Error())
		return result, err
	}
	return result, err
}
