package util

import (
	"aDi/log"
	jsoniter "github.com/json-iterator/go"
)

// MarshalToStringWithOutErr 不返回错误的marshal
func MarshalToStringWithOutErr(v interface{}) string {
	info, err := jsoniter.MarshalToString(v)
	if err != nil {
		log.Errorf("marshal json fail,err:%s", err.Error())
		return info
	}
	return info
}

// MarshalWithoutErr 不返回错误的marshal
func MarshalWithoutErr(v interface{}) []byte {
	info, err := jsoniter.Marshal(v)
	if err != nil {
		log.Errorf("marshal fail,err:%s", err.Error())
		return info
	}
	return info
}
