package model

import "net/http"

// ErrCode 错误码
type ErrCode int32

// ErrMsg 错误信息
type ErrMsg string

type CodeMsg struct {
	Code int32  `json:"code"` // 返回码
	Msg  string `json:"msg"`  // 返回消息
}

const (
	// err code enum
	ErrCodeSuccess = ErrCode(0)
)

var (
	ErrICodeSuccess  = CodeMsg{Code: 0, Msg: "success"}        // 成功
	ErrIInner        = CodeMsg{Code: 4, Msg: "内部错误"}           // 内部错误
	ErrIParse        = CodeMsg{Code: 8, Msg: "解析失败"}           // 解析失败
	ErrIInvalidParam = CodeMsg{Code: 9, Msg: "无效参数"}           // 无效参数
	ErrINoDb         = CodeMsg{Code: 10, Msg: "没有对应的database"} // 没有database
	ErrIDbFail       = CodeMsg{Code: 12, Msg: "数据库操作失败"}       // db操作失败
	ErrIS2S          = CodeMsg{Code: 13, Msg: "调用其他服务失败"}      // 调用其他服务失败
	ErrINoAuth       = CodeMsg{Code: 21, Msg: "没有权限"}          // 没有权限
	ErrIGetAuthFail  = CodeMsg{Code: 22, Msg: "获取权限失败"}        // 获取权限失败
	ErrINotFound     = CodeMsg{Code: http.StatusNotFound, Msg: "404 page not found"}
	ErrINoLogin      = CodeMsg{Code: 108, Msg: "未进行登录，请先登录"} // 注册前需要先进性登录
	ErrInvalidToken  = CodeMsg{Code: 109, Msg: "无效token"}
)
