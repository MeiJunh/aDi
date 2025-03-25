package model

// Conversation 会话信息
type Conversation struct {
	Id              int64           `json:"id"`              // 会话id
	LastMsg         string          `json:"lastMsg"`         // 最后一条记录
	UnreadNum       int32           `json:"unreadNum"`       // 未读数量
	ChatLeftNum     int64           `json:"chatLeftNum"`     // 剩余聊天次数
	DigitalChatInfo ChatDigitalInfo `json:"digitalChatInfo"` // 数字人信息
}

// ChatDigitalInfo 聊天时数字人的简要信息
type ChatDigitalInfo struct {
	Uid         int64  `json:"uid"`  // 数字人对应的用户id
	Name        string `json:"name"` // 数字人名字
	Icon        string `json:"icon"` // 头像
	ConChatConf        // 所选择的对话套餐配置
}

// DigitalConversation -- 数字人视角的会话信息
type DigitalConversation struct {
	Id               int64         `json:"id"`
	CharUserBaseInfo *UserBaseInfo `json:"charUserBaseInfo"` // 聊天人的基础信息
	IsAnonymity      Switch        `json:"isAnonymity"`      // 是否是匿名
}

// GetMessageListReq 获取自己某次会话的聊天记录
type GetMessageListReq struct {
	ConversationId int64  `json:"conversationId"` // 会话id -- 按用户的会话id指的是会话配置id，按机器人查看的是真是会话表中的id
	Index          string `json:"index"`          // 游标透传 -- 首次为空字符串
	PageSize       int    `json:"pageSize"`       // 数量
}

// ChatMessage 聊天信息 -- 一轮问答
type ChatMessage struct {
	Id         int64    `json:"id"`                // 对应的一轮问答id
	TagList    []string `json:"tagList,omitempty"` // 标签列表
	UMessage   string   `json:"uMessage"`          // 用户说的
	DMessage   string   `json:"dMessage"`          // 数字人说的
	CreateTime int64    `json:"createTime"`        // 问答创建时间
}

// AiChatGenerateReq ai问答生成
type AiChatGenerateReq struct {
	DigitalUid int64  `json:"digitalUid"` // 数字人对应的uid
	Content    string `json:"content"`    // 问话内容
}

// AiChatGenerateRsp ai问答生成返回
type AiChatGenerateRsp struct {
	KeyCode   string `json:"keyCode"`   // 异步获取结果使用的key
	MessageId int64  `json:"messageId"` // 当前问答对应的id
}

type AiAsyncResult struct {
	AiStatus ErrCode `json:"aiStatus"`           // ai识别结果 -- 是否需要继续查询
	AiErrMsg string  `json:"aiErrMsg,omitempty"` // ai识别错误结果显示
	AiResult string  `json:"aiResult,omitempty"` // ai识别结果
}

// DbConversationConf 当前会话配置信息 -- 用作个人与机器人聊天会话列表
type DbConversationConf struct {
	ID             int64        `json:"id" db:"id"`                       // 会话配置id -- 作用与个人的会话id使用
	UID            int64        `json:"uid" db:"uid"`                     // 聊天人uid
	DigitalUID     int64        `json:"digitalUid" db:"digital_uid"`      // 机器人uid
	IsAnonymity    Switch       `json:"isAnonymity" db:"is_anonymity"`    // 是否匿名,0不匿名,1匿名
	ChatConf       string       `json:"chatConf" db:"chat_conf"`          // 聊天配置
	LastMsg        string       `json:"lastMsg" db:"last_msg"`            // 最后一条聊天记录
	ChatTotalNum   int64        `json:"chatTotalNum" db:"chat_total_num"` // 聊天总次数
	ChatUseNum     int64        `json:"chatUseNum" db:"chat_use_num"`     // 聊天已用次数
	ChatConfStruct *ConChatConf `json:"chatConfStruct"`                   // 聊天会话配置
}

// ConChatConf 会话聊天设置
type ConChatConf struct {
	KindList    []string `json:"kindList"`    // 所选择的对话套餐
	Scene       string   `json:"scene"`       // 场景
	Style       string   `json:"style"`       // 风格
	Appellation string   `json:"appellation"` // 称呼
}

// DbMessage 消息信息
type DbMessage struct {
	ID             int64  `json:"id" db:"id"`                          // 消息id
	UID            int64  `json:"uid" db:"uid"`                        // 聊天人uid
	DigitalUID     int64  `json:"digitalUid" db:"digital_uid"`         // 机器人uid
	ConversationID int64  `json:"conversationId" db:"conversation_id"` // 所属会话的ID
	UMessage       string `json:"uMessage" db:"u_message"`             // 用户说的
	DMessage       string `json:"dMessage" db:"d_message"`             // 机器人说的
	ChatConf       string `json:"chatConf" db:"chat_conf"`             // 聊天配置
	VoiceUrl       string `json:"voiceUrl" db:"voice_url"`             // 语音条
	CreateTime     int64  `json:"createTime" db:"create_time"`         // 时间
}

// DbConversation 会话列表
type DbConversation struct {
	ID          int64  `json:"id" db:"id"`                    // 会话id
	UID         int64  `json:"uid" db:"uid"`                  // 聊天人uid
	DigitalUID  int64  `json:"digitalUid" db:"digital_uid"`   // 机器人uid
	IsAnonymity Switch `json:"isAnonymity" db:"is_anonymity"` // 是否匿名,0不匿名,1匿名
	CreateTime  int64  `json:"createTime" db:"create_time"`   // 会话创建时间
}

// BuyChatNumReq 购买聊天次数入参
type BuyChatNumReq struct {
	DigitalUid     int64        `json:"digitalUid"`     // 需要聊天的机器人id
	IsAnonymity    Switch       `json:"isAnonymity"`    // 是否匿名,0不匿名,1匿名
	ChatConfStruct *ConChatConf `json:"chatConfStruct"` // 聊天会话配置
	BuyNum         int64        `json:"buyNum"`         // 购买数量
	TotalPrice     int64        `json:"totalPrice"`     // 总价格 -- 会进行校验
}

type GetChatVoiceRsp struct {
	VoiceUrl string `json:"voiceUrl"` // 语音条
}
