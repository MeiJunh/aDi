package util

import (
	"aDi/log"
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// ToInt64 转化成int64
func ToInt64(s string) (i int64) {
	var err error
	i, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Errorf("parse fail,s:%s,err:%s", s, err.Error())
		return i
	}
	return i
}

// JoinInt64 合并int64
func JoinInt64(list []int64, sep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(list[0], 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatInt(list[0], 10))
	for _, s := range list[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.FormatInt(s, 10))
	}
	return b.String()
}

// CopyGetRequestBody 拷贝
func CopyGetRequestBody(request *http.Request) []byte {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil
	}

	_ = request.Body.Close()
	request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}
