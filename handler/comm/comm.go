package comm

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// ReadBodyFromGin 从gin ctx中获取body并解析
func ReadBodyFromGin(c *gin.Context, req interface{}) (bodyBuff []byte, err error) {
	// 获取body info
	bodyBuff, err = c.GetRawData()
	if err != nil {
		return bodyBuff, err
	}

	// 解析
	err = jsoniter.Unmarshal(bodyBuff, req)
	return bodyBuff, err
}
