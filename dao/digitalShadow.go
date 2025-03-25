package dao

import (
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// GetDigitalInfo 获取数字人信息
func GetDigitalInfo(uid int64) (info *model.DbDigitalInfo, err error) {
	info = &model.DbDigitalInfo{
		ChargeConfInfo: &model.ChargeInfoConf{},
		DAllKinds:      &model.DigitalAllKinds{},
	}
	err = dbClient.FindOneWithNull(info, "SELECT id, uid,icon, digital_name, can_anonymity, prologue,clone_voice, charge_conf, digital_all_kinds,status"+
		" FROM t_digital_info where uid = ?", uid)
	if err != nil {
		log.Errorf("select fail, err:%s", err.Error())
		return info, err
	}
	if info.ChargeConf != "" {
		_ = jsoniter.UnmarshalFromString(info.ChargeConf, info.ChargeConfInfo)
	}
	if info.DigitalAllKinds != "" {
		_ = jsoniter.UnmarshalFromString(info.DigitalAllKinds, info.DAllKinds)
	}
	return info, err
}

// UpsertDigitalInfo 新增/更新数字人信息
func UpsertDigitalInfo(uid int64, info *model.DigitalShadowInfo) (err error) {
	// 先判断数字人信息是否存在
	rInfo, err := GetDigitalInfo(uid)
	if err != nil {
		log.Errorf("get digital info fail, err:%s", err.Error())
		return err
	}
	if rInfo == nil || rInfo.Uid <= 0 {
		// 走新增逻辑
		err = AddDigitalInfo(uid, info)
		return err
	}
	// 走更新逻辑
	err = UpdateDigitalInfo(uid, info)
	return err
}

// AddDigitalInfo 新增数字人信息
func AddDigitalInfo(uid int64, info *model.DigitalShadowInfo) (err error) {
	_, err = dbClient.Exec("INSERT INTO t_digital_info (uid,icon, digital_name, can_anonymity, prologue, charge_conf,digital_all_kinds,status"+
		") VALUES (?,?, ?, ?, ?, ?, ?, ?)", uid, info.UserInfo.GetIcon(), info.DigitalName, info.CanAnonymity, info.Prologue, util.MarshalToStringWithOutErr(info.ChargeConf),
		util.MarshalToStringWithOutErr(info.DigitalAllKinds), model.DSCreate)
	if err != nil {
		log.Errorf("add digital info fail, err:%s", err.Error())
		return err
	}
	return nil
}

// UpdateDigitalInfo 更新数字人所有信息
func UpdateDigitalInfo(uid int64, info *model.DigitalShadowInfo) (err error) {
	_, err = dbClient.Exec("update t_digital_info set digital_name = ?, can_anonymity = ?, prologue = ?, charge_conf = ?,"+
		"digital_all_kinds = ?,status = ? where uid = ?", info.DigitalName, info.CanAnonymity, info.Prologue, util.MarshalToStringWithOutErr(info.ChargeConf),
		util.MarshalToStringWithOutErr(info.DigitalAllKinds), model.DSCreate, uid)
	if err != nil {
		log.Errorf("update digital info fail, err:%s", err.Error())
		return err
	}
	return err
}

type DigitalField string

const (
	DFDigitalName     = DigitalField("digital_name")
	DFCanAnonymity    = DigitalField("can_anonymity")
	DFPrologue        = DigitalField("prologue")
	DFCloneVoice      = DigitalField("clone_voice")
	DFChargeConf      = DigitalField("charge_conf")
	DFDigitalAllKinds = DigitalField("digital_all_kinds")
)

// UpdateDigitalField 更新数字人指定字段信息 value 只能为int或者string
func UpdateDigitalField(uid int64, field DigitalField, value interface{}) (err error) {
	_, err = dbClient.Exec(fmt.Sprintf("update t_digital_info set %s = ? where uid = ?", field), value, uid)
	if err != nil {
		log.Errorf("update digital field fail, err:%s", err.Error())
		return err
	}
	return err
}

// GetDigitalListByUid 根据uid列表获取数字人信息
func GetDigitalListByUid(uidList []int64) (list []*model.DbDigitalInfo, err error) {
	if len(uidList) <= 0 {
		return list, err
	}
	// 获取数字人信息列表
	list, err = getDigitalList(fmt.Sprintf(" where uid in (%s)", util.JoinInt64(uidList, ",")))
	if err != nil {
		log.Errorf("get digital list fail, err:%s", err.Error())
		return list, err
	}
	return list, err
}

// getDigitalList 获取数字人信息列表
func getDigitalList(whereQuery string, params ...interface{}) (list []*model.DbDigitalInfo, err error) {
	list = make([]*model.DbDigitalInfo, 0)
	err = dbClient.FindList(&list, "SELECT id, uid,icon, digital_name, can_anonymity, prologue,clone_voice, charge_conf, digital_all_kinds,status"+
		" FROM t_digital_info "+whereQuery, params...)
	if err != nil {
		log.Errorf("get list fail, err:%s", err.Error())
		return list, err
	}
	return list, err
}
