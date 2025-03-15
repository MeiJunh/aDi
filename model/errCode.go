package model

import "net/http"

// ErrCode 错误码
type ErrCode int32

// ErrMsg 错误信息
type ErrMsg string

type CodeMsg struct {
	Code ErrCode `json:"code"` // 返回码
	Msg  string  `json:"msg"`  // 返回消息
}

const (
	// err code enum
	ErrCodeSuccess  = ErrCode(0)
	ErrCodeSuccess2 = ErrCode(1)
	ECInner         = ErrCode(4)
	ECParse         = ErrCode(8)
	ECInvalidParam  = ErrCode(9)
	ECDbFail        = ErrCode(12)
	ECS2S           = ErrCode(13)
	ECBan           = ErrCode(14) // 操作被禁止
	ECNoExist       = ErrCode(15) // 对象不存在
	ECHttpDo        = ErrCode(16) // http请求失败
	ECNoAuth        = ErrCode(21)
	ECNoLogin       = ErrCode(108)
	ECInvalidToken  = ErrCode(109)
)

var (
	ErrICodeSuccess  = CodeMsg{Code: ErrCodeSuccess, Msg: "success"} // 成功
	ErrIInner        = CodeMsg{Code: ECInner, Msg: "内部错误"}           // 内部错误
	ErrIParse        = CodeMsg{Code: ECParse, Msg: "解析失败"}           // 解析失败
	ErrIInvalidParam = CodeMsg{Code: ECInvalidParam, Msg: "无效参数"}    // 无效参数
	ErrIDbFail       = CodeMsg{Code: ECDbFail, Msg: "数据库操作失败"}       // db操作失败
	ErrIS2S          = CodeMsg{Code: ECS2S, Msg: "调用其他服务失败"}         // 调用其他服务失败
	ErrINoAuth       = CodeMsg{Code: ECNoAuth, Msg: "没有权限"}          // 没有权限
	ErrIGetAuthFail  = CodeMsg{Code: 22, Msg: "获取权限失败"}              // 获取权限失败
	ErrINotFound     = CodeMsg{Code: http.StatusNotFound, Msg: "404 page not found"}
	ErrINoLogin      = CodeMsg{Code: ECNoLogin, Msg: "未进行登录，请先登录"} // 注册前需要先进性登录
	ErrInvalidToken  = CodeMsg{Code: ECInvalidToken, Msg: "无效token"}
)
