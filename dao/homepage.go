package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
)

// UpsertHomepageConf 新增或者更新主页信息
func UpsertHomepageConf(uid int64, conf *model.HpConf) (err error) {
	// 先获取用户主页配置
	hpConf, err := GetHomepageConfByUid(uid)
	if err != nil {
		log.Error("get homepage conf fail,err:%s", err.Error())
		return err
	}
	// 如果不存在则进行新增
	if hpConf == nil || hpConf.ID <= 0 {
		err = AddHomepageConf(uid, conf)
		if err != nil {
			log.Error("add homepage conf fail,err:%s", err.Error())
			return err
		}
		return err
	}
	// 否则进行更新
	err = UpdateHomepageConf(uid, conf)
	if err != nil {
		log.Error("update homepage conf fail,err:%s", err.Error())
		return err
	}
	return err
}

// AddHomepageConf 新增用户主页配置
func AddHomepageConf(uid int64, conf *model.HpConf) (err error) {
	_, err = dbClient.Exec("insert into  t_homepage_conf (uid, conf_str) values (?, ?)", uid, util.MarshalToStringWithOutErr(conf))
	if err != nil {
		log.Errorf("add homepage conf fail,err: %s", err.Error())
		return err
	}
	return err
}

// UpdateHomepageConf 编辑用户主页配置
func UpdateHomepageConf(uid int64, confStr *model.HpConf) (err error) {
	_, err = dbClient.Exec("update t_homepage_conf set conf_str = ? where uid = ?", util.MarshalToStringWithOutErr(confStr), uid)
	if err != nil {
		log.Errorf("update homepage conf fail,err: %s", err.Error())
		return err
	}
	return err
}

// GetHomepageConfByUid 获取指定用户的主页配置
func GetHomepageConfByUid(uid int64) (info *model.DbHomepageConf, err error) {

	return info, err
}
