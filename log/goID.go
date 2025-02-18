package log

import (
	"bytes"
	"github.com/v2pro/plz/gls"
	"strconv"
	"strings"
)

const (
	logCommonKeyGoID = "GoID"
)

// GetGoID 测试对比，这种方式性能会好一点
func GetGoID() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("GoID:")
	buf.WriteString(strconv.FormatInt(gls.GoID(), 10))
	return buf.String()
}

func GetFuncName(name string) string {
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx:]
		idx = strings.IndexByte(name, '.')
		if idx != -1 {
			name = strings.TrimPrefix(name[idx:], ".")
		}
	}
	return name
}
