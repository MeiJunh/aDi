package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"fmt"
	"strconv"
)

// AddFollowInfo 关注
func AddFollowInfo(uid, follower int64) {
	result, err := dbClient.Exec("insert into t_follow(uid,follower) values(?,?)", uid, follower)
	if err != nil {
		log.Errorf("add follow info fail,err: %s", err.Error())
		return
	}
	effect, err := result.RowsAffected()
	if err != nil {
		log.Errorf("add follow info fail,err: %s", err.Error())
		return
	}
	if effect > 0 {
		// 添加关注统计
		UpdateSocialStatisticInfo(uid, 1, FieldFollow)
	}
	return
}

// CancelFollow 取消关注
func CancelFollow(uid, follower int64) {
	result, err := dbClient.Exec("delete from t_follow where uid = ? and follower = ?", uid, follower)
	if err != nil {
		log.Errorf("cancel follow fail,err: %s", err.Error())
		return
	}
	effect, err := result.RowsAffected()
	if err != nil {
		log.Errorf("cancel follow info fail,err: %s", err.Error())
		return
	}
	if effect > 0 {
		// 取消关注统计
		UpdateSocialStatisticInfo(uid, -1, FieldFollow)
	}
	return
}

// AddFavorInfo 点赞
func AddFavorInfo(uid, liker int64) {
	result, err := dbClient.Exec("insert into t_favor(uid,liker) values(?,?)", uid, liker)
	if err != nil {
		log.Errorf("add favor info fail,err: %s", err.Error())
		return
	}
	effect, err := result.RowsAffected()
	if err != nil {
		log.Errorf("add favor info fail,err: %s", err.Error())
		return
	}
	if effect > 0 {
		// 添加点赞统计
		UpdateSocialStatisticInfo(uid, 1, FieldFavor)
	}
	return
}

// AddViewRecord 添加浏览记录
func AddViewRecord(uid, viewer int64) {
	result, err := dbClient.Exec("insert into t_view_record(uid,viewer) values(?,?)", uid, viewer)
	if err != nil {
		log.Errorf("add view info fail,err: %s", err.Error())
		return
	}
	effect, err := result.RowsAffected()
	if err != nil {
		log.Errorf("add view info fail,err: %s", err.Error())
		return
	}
	if effect > 0 {
		// 添加浏览统计
		UpdateSocialStatisticInfo(uid, 1, FieldView)
	}
	return
}

// GetMyFollowList 获取我的关注列表，follower是我
func GetMyFollowList(uid int64, index string, pageSize int) (list []*model.DbFollow, hasMore bool, nextIndex string, err error) {
	if pageSize <= 0 || pageSize > 100 {
		// 进行格式化
		pageSize = 20
	}

	// 当前游标用作offset
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	list, err = getFollowList("where follower = ? order by create_time desc limit ?,?", uid, offset, pageSize)
	if err != nil {
		log.Errorf("get follow list fail,err: %s", err.Error())
		return list, hasMore, nextIndex, err
	}

	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.Itoa(int(offset))
	}
	return list, hasMore, nextIndex, nil
}

// GetFollowMeList 获取关注我的列表
func GetFollowMeList(uid int64, index string, pageSize int) (list []*model.DbFollow, hasMore bool, nextIndex string, err error) {
	if pageSize <= 0 || pageSize > 100 {
		// 进行格式化
		pageSize = 20
	}

	// 当前游标用作offset
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	list, err = getFollowList("where uid = ? order by create_time desc limit ?,?", uid, offset, pageSize)
	if err != nil {
		log.Errorf("get follow list fail,err: %s", err.Error())
		return list, hasMore, nextIndex, err
	}

	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.Itoa(int(offset))
	}
	return list, hasMore, nextIndex, nil
}

// getFollowList 获取关注列表
func getFollowList(whereQuery string, param ...interface{}) (list []*model.DbFollow, err error) {
	list = make([]*model.DbFollow, 0)
	err = dbClient.FindList(&list, "select uid,follower,UNIX_TIMESTAMP(create_time) as ctime from t_follow "+whereQuery, param...)
	if err != nil {
		log.Errorf("get follow list fail,err: %s", err.Error())
		return list, err
	}
	return list, err
}

// GetFavorMeList 获取点赞我的列表
func GetFavorMeList(uid int64, index string, pageSize int) (list []*model.DbFavor, hasMore bool, nextIndex string, err error) {
	if pageSize <= 0 || pageSize > 100 {
		// 进行格式化
		pageSize = 20
	}

	// 当前游标用作offset
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	list = make([]*model.DbFavor, 0)
	err = dbClient.FindList(&list, "select uid,liker,UNIX_TIMESTAMP(create_time) as ctime from t_favor where uid = ? "+
		"order by create_time desc limit ?,?", uid, offset, pageSize)
	if err != nil {
		log.Errorf("get favor list fail,err: %s", err.Error())
		return list, hasMore, nextIndex, err
	}

	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.Itoa(int(offset))
	}
	return list, hasMore, nextIndex, nil
}

// GetViewMeList 获取浏览我的列表
func GetViewMeList(uid int64, index string, pageSize int) (list []*model.DbView, hasMore bool, nextIndex string, err error) {
	if pageSize <= 0 || pageSize > 100 {
		// 进行格式化
		pageSize = 20
	}

	// 当前游标用作offset
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	list, err = getViewList(" where uid = ? order by create_time limit ?,?", uid, offset, pageSize)
	if err != nil {
		log.Errorf("get view list fail,err: %s", err.Error())
		return list, hasMore, nextIndex, err
	}

	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.Itoa(int(offset))
	}
	return list, hasMore, nextIndex, nil
}

// GetMyViewList 获取我浏览的列表
func GetMyViewList(uid int64, index string, pageSize int) (list []*model.DbView, hasMore bool, nextIndex string, err error) {
	if pageSize <= 0 || pageSize > 100 {
		// 进行格式化
		pageSize = 20
	}

	// 当前游标用作offset
	offset := util.ToInt64(index)
	if offset < 0 {
		offset = 0
	}
	list, err = getViewList(" where viewer = ? order by create_time limit ?,?", uid, offset, pageSize)
	if err != nil {
		log.Errorf("get my view list fail,err: %s", err.Error())
		return list, hasMore, nextIndex, err
	}

	if len(list) >= pageSize {
		hasMore = true
		nextIndex = strconv.Itoa(int(offset))
	}
	return list, hasMore, nextIndex, nil
}

func getViewList(whereQuery string, param ...interface{}) (list []*model.DbView, err error) {
	list = make([]*model.DbView, 0)
	err = dbClient.FindList(&list, "select uid,viewer,UNIX_TIMESTAMP(create_time) as ctime from t_view_record "+whereQuery, param...)
	if err != nil {
		log.Errorf("get view list fail,err: %s", err.Error())
		return list, err
	}
	return list, err
}

// GetSocialStatisticInfo 获取统计信息
func GetSocialStatisticInfo(uid int64) (info *model.DBFollowStatisticInfo, err error) {
	info = &model.DBFollowStatisticInfo{}
	err = dbClient.FindOneWithNull(info, "select id,uid,follow_num,favor_num,view_num from t_social_statistic where uid = ?", uid)
	if err != nil {
		log.Errorf("get follow statistic info fail,err: %s", err.Error())
		return info, err
	}
	return info, err
}

// InitSocialStatisticInfo 初始化统计信息
func InitSocialStatisticInfo(uid int64) {
	_, err := dbClient.Exec("insert into t_social_statistic (uid) value (?)", uid)
	if err != nil {
		log.Errorf("init follow statistic info fail,err: %s", err.Error())
		return
	}
	return
}

// SocialStatisticField 统计表更新字段
type SocialStatisticField string

const (
	FieldFollow = SocialStatisticField("follow_num")
	FieldFavor  = SocialStatisticField("favor_num")
	FieldView   = SocialStatisticField("view_num")
)

// UpdateSocialStatisticInfo 更新社交统计信息
func UpdateSocialStatisticInfo(uid int64, num int64, field SocialStatisticField) {
	_, err := dbClient.Exec(fmt.Sprintf("update t_social_statistic set %s = %s + ? where uid = ?", field, field), num, uid)
	if err != nil {
		log.Errorf("update follow statistic info fail,err: %s", err.Error())
		return
	}
	return
}
