package comm

import (
	"aDi/log"
	"aDi/service"
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

// GetUidFromCon 从context 中获取uid
func GetUidFromCon(c *gin.Context) (uid int64) {
	var err error
	// 从token中解析有效的uid
	uid, err = service.ValidateToken(c.GetHeader("token"))
	if err != nil {
		log.Errorf("get uid from con fail,token:%s,err:%s", c.GetHeader("token"), err.Error())
		return uid
	}
	return uid
}
