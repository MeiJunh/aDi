package model

// BaseRsp 公共出参
type BaseRsp struct {
	CodeMsg
	Data interface{} `json:"data"` // 返回body
}

// WriteMsg 写数据 -- 出错时不需要写data返回值
func (r *BaseRsp) WriteMsg(errInfo CodeMsg) {
	r.CodeMsg = errInfo
}

// Generate 赋值并且返回
func (r *BaseRsp) Generate(errInfo CodeMsg) *BaseRsp {
	r.WriteMsg(errInfo)
	return r
}

// Switch 开关
type Switch int32

const (
	SwitchOff = Switch(0) // 关闭
	SwitchOn  = Switch(1) // 打开
)
