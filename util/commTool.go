package util

import (
	"aDi/log"
	jsoniter "github.com/json-iterator/go"
	"time"
)

// 通用的工具方法

// TimeCost 函数耗时以及出入参打印,params为附带参数,比如input中只想打印部分,不想打印其中的数组,则可以在input中不填
func TimeCost(funcName string, input, output interface{}, params ...string) func() {
	now := time.Now()
	return func() {
		inputStr, _ := jsoniter.MarshalToString(input)
		outputStr, _ := jsoniter.MarshalToString(output)
		if len(params) == 0 {
			log.Debugf("%s,cost:%d ms,input:%s,output:%s", funcName, time.Since(now).Milliseconds(), inputStr, outputStr)
		} else {
			log.Debugf("%s,cost:%d ms,input:%s,output:%s,params:%v", funcName, time.Since(now).Milliseconds(), inputStr, outputStr, params)
		}
	}
}
