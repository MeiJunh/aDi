package model

// GameState 游戏状态
type GameState int32

// GameChatState 游戏聊天状态 -- 表示当前回合结果
type GameChatState int32

const (
	GSDefault  = GameState(0)     // 正常
	GSDel      = GameState(-1)    // 被删除
	GSNoRE     = GameState(1)     // 红包已经发完
	GSExpire   = GameState(2)     // 游戏过期
	GSPlayOver = GameState(10)    // 当前用户已经将该游戏玩完了
	GSPlayWin  = GameState(11)    // 当前用户已经获得了该游戏的奖励
	GCSWrong   = GameChatState(1) // 错误
	GCSRight   = GameChatState(2) // 正确
)

// GameInfo 游戏相关信息
type GameInfo struct {
	Id            int64     `json:"id"`            // 游戏id
	Name          string    `json:"name"`          // 游戏名
	RETotalAmount int64     `json:"RETotalAmount"` // 红包总金额 -- 单位分
	RETotalNum    int64     `json:"RETotalNum"`    // 红包个数
	REClaimNum    int64     `json:"REClaimNum"`    // 红包被领取的数量
	Prologue      string    `json:"prologue"`      // 开场白
	AnswerList    []string  `json:"answerList"`    // 答案列表
	State         GameState `json:"state"`         // 游戏状态
	CreateTime    int64     `json:"createTime"`    // 创建时间
}

// GamePlayInfo 玩游戏的相关信息
type GamePlayInfo struct {
	GameName    string           `json:"gameName"`    // 游戏名
	Prologue    string           `json:"prologue"`    // 开场白
	ChatList    []*GameChatInfo  `json:"chatList"`    // 游戏的聊天信息
	GameState   GameState        `json:"gameState"`   // 当前游戏状态
	DigitalInfo *ChatDigitalInfo `json:"digitalInfo"` // 数字人信息
}

// GameChatInfo 游戏聊天信息
type GameChatInfo struct {
	Input       string        `json:"input"`
	Output      string        `json:"output"`
	ResultState GameChatState `json:"resultState"` // 回答结果
	CreateTime  int64         `json:"createTime"`
}

// DbGame 游戏配置表
type DbGame struct {
	Id            int64     `json:"id" db:"id"`                         //
	Uid           int64     `json:"uid" db:"uid"`                       // uid
	Name          string    `json:"name" db:"name"`                     // 游戏名
	Prologue      string    `json:"prologue" db:"prologue"`             // 开场白
	ReTotalAmount int64     `json:"reTotalAmount" db:"re_total_amount"` // 红包总金额 -- 单位分
	ReTotalNum    int64     `json:"reTotalNum" db:"re_total_num"`       // 红包总数量
	ReClaimNum    int64     `json:"reClaimNum" db:"re_claim_num"`       // 红包被领取的数量
	AnswerListStr string    `json:"answerListStr" db:"answer_list_str"` // 答案信息
	State         GameState `json:"state" db:"state"`                   // 游戏状态
	Version       int64     `json:"version" db:"version"`               // 版本号
	CreateTime    int64     `json:"createTime" db:"create_time"`        //
	AnswerList    []string  `json:"-"`
}

func TransDbGameToGameInfo(list []*DbGame) (rList []*GameInfo) {
	rList = make([]*GameInfo, 0, len(list))
	for _, v := range list {
		rList = append(rList, &GameInfo{
			Id:            v.Id,
			Name:          v.Name,
			RETotalAmount: v.ReTotalAmount,
			RETotalNum:    v.ReTotalNum,
			REClaimNum:    v.ReClaimNum,
			Prologue:      v.Prologue,
			AnswerList:    v.AnswerList,
			State:         v.State,
			CreateTime:    v.CreateTime,
		})
	}
	return rList
}

// ChatWithGameReq 和游戏聊天入参
type ChatWithGameReq struct {
	ShareCode string `json:"shareCode"` // 分享码
	Input     string `json:"input"`     // 回答
}

// ChatWithGameRsp 和游戏聊天回参
type ChatWithGameRsp struct {
	Output      string        `json:"output"`
	ResultState GameChatState `json:"resultState"` // 回答结果
	REAmount    int64         `json:"reAmount"`    // 红包金额
	GameState   GameState     `json:"gameState"`   // 游戏状态
}

// DbGamePlayRecord 游戏轮次记录表
type DbGamePlayRecord struct {
	ID          int64         `json:"id" db:"id"`                    //
	Uid         int64         `json:"uid" db:"uid"`                  // uid
	GameID      int64         `json:"gameId" db:"game_id"`           // 游戏id
	Input       string        `json:"input" db:"input"`              // 输入 -- 用户的回答
	Output      string        `json:"output" db:"output"`            // 输出结果
	ResultState GameChatState `json:"resultState" db:"result_state"` // 游戏结果1:错误,2:正确
	CreateTime  int64         `json:"createTime" db:"create_time"`   //
}
