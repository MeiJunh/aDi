package model

import jsoniter "github.com/json-iterator/go"

// DigitalStatus 数字分身状态
type DigitalStatus int32

const (
	DSNoCreate = DigitalStatus(0) // 未创建
	DSCreate   = DigitalStatus(1) // 已创建
)

// DigitalShadowInfo 数字分身所有信息
type DigitalShadowInfo struct {
	Uid             int64            `json:"uid"`
	UserInfo        *UserAllInfo     `json:"userInfo"`        // 用户信息
	CloneVoice      string           `json:"cloneVoice"`      // 克隆语音条
	Prologue        string           `json:"prologue"`        // 开场白
	DigitalName     string           `json:"digitalName"`     // 分身名
	CanAnonymity    Switch           `json:"canAnonymity"`    // 是否可以匿名
	ChargeConf      *ChargeInfoConf  `json:"chargeConf"`      // 收费设置
	DigitalAllKinds *DigitalAllKinds `json:"digitalAllKinds"` // 百变分身配置
	Status          DigitalStatus    `json:"status"`          // 0 未创建，1 已创建
}

// DigitalAllKinds 百变分身配置
type DigitalAllKinds struct {
	SceneList       []string `json:"sceneList"`       // 场景
	StyleList       []string `json:"styleList"`       // 风格
	AppellationList []string `json:"appellationList"` // 称呼
}

// ChargeInfoConf 收费折信息设置
type ChargeInfoConf struct {
	BasePrice       int64         `json:"basePrice"`       // 基础价格
	SceneList       []*ChargeItem `json:"sceneList"`       // 场景
	StyleList       []*ChargeItem `json:"styleList"`       // 风格
	AppellationList []*ChargeItem `json:"appellationList"` // 称呼
	CanAnonymity    Switch        `json:"canAnonymity"`    // 是否可以匿名
	AnonymityPrice  int64         `json:"anonymityPrice"`  // 匿名价格
}

type ChargeItem struct {
	ItemName string `json:"itemName"` // 对应的收费名
	Price    int64  `json:"price"`    // 单价 -- 分
}

// DbDigitalInfo 数据库中的数字人信息
type DbDigitalInfo struct {
	ID              int64            `json:"id" db:"id"`                             // 主键ID
	Uid             int64            `json:"uid" db:"uid"`                           // 数字人uid
	Icon            string           `json:"icon" db:"icon"`                         // 头像
	DigitalName     string           `json:"digitalName" db:"digital_name"`          // 分身名
	CanAnonymity    Switch           `json:"canAnonymity" db:"can_anonymity"`        // 是否可以匿名
	Prologue        string           `json:"prologue" db:"prologue"`                 // 开场白
	CloneVoice      string           `json:"cloneVoice" db:"clone_voice"`            // 克隆语音条
	ChargeConf      string           `json:"chargeConf" db:"charge_conf"`            // 收费设置
	DigitalAllKinds string           `json:"digitalAllKinds" db:"digital_all_kinds"` // 百变分身设置
	Status          DigitalStatus    `json:"status" db:"status"`                     // 0 未创建，1 已创建
	ChargeConfInfo  *ChargeInfoConf  `json:"-"`                                      // 配置信息
	DAllKinds       *DigitalAllKinds `json:"-"`                                      // 所有的配置信息
}

func (d *DbDigitalInfo) GetDigitalName() string {
	if d == nil {
		return ""
	}
	return d.DigitalName
}

func (d *DbDigitalInfo) GetIcon() string {
	if d == nil {
		return ""
	}
	return d.Icon
}

// DbDigitalInfoTrans 数字人信息转化--不能返回nil
func DbDigitalInfoTrans(info *DbDigitalInfo) (r *DigitalShadowInfo) {
	if info == nil {
		return &DigitalShadowInfo{}
	}
	r = &DigitalShadowInfo{
		Uid:             info.Uid,
		CloneVoice:      info.CloneVoice,
		Prologue:        info.Prologue,
		DigitalName:     info.DigitalName,
		CanAnonymity:    info.CanAnonymity,
		ChargeConf:      &ChargeInfoConf{},
		DigitalAllKinds: &DigitalAllKinds{},
		Status:          0,
	}
	if info.ChargeConf != "" {
		_ = jsoniter.UnmarshalFromString(info.ChargeConf, r.ChargeConf)
	}
	if info.DigitalAllKinds != "" {
		_ = jsoniter.UnmarshalFromString(info.DigitalAllKinds, r.DigitalAllKinds)
	}
	return r
}
