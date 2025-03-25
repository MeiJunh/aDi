package noLoginApi

import (
	"aDi/model"
	"github.com/gin-gonic/gin"
)

/*
一些通用下拉列表
*/

// GetProvinceCityMap 获取省市下拉列表
func (n *NoLoginHandlerImp) GetProvinceCityMap(c *gin.Context) (rsp model.BaseRsp) {
	rsp.Data = map[string][]string{}
	return
}
