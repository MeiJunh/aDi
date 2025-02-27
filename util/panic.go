package util

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var (
	mGoPath    string
	mGoRoot    string
	mGoPathLen int
	mGoRootLen int
)

func init() {
	mGoPath = os.Getenv("GOPATH")
	mGoRoot = os.Getenv("GOROOT")
	mGoPathLen = len(mGoPath)
	mGoRootLen = len(mGoRoot)
}

// CallerInfo 根据pc获取调用信息
func CallerInfo(pc uintptr) string {
	var buf bytes.Buffer

	fun := runtime.FuncForPC(pc)
	if nil == fun {
		return ""
	}
	fn := fun.Name()
	if strings.HasPrefix(fn, "runtime.") {
		return ""
	}

	f, l := fun.FileLine(pc) // pc保存的是下一个地址，所以要回退
	if strings.HasPrefix(f, mGoPath) {
		f = f[mGoPathLen+1:]
	} else if strings.HasPrefix(f, mGoRoot) {
		f = f[mGoRootLen+1:]
	}

	buf.WriteString(fn)
	buf.WriteString("(")
	buf.WriteString(f)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(l))
	buf.WriteString(")")

	return buf.String()
}

// CallStack 生成可以用于log的调用栈信息（没有换行）
func CallStack(skip, depth int) string {
	if depth < 2 {
		depth = 2
	}
	if skip < 0 {
		skip = 0
	}

	var buf bytes.Buffer
	fpcs := make([]uintptr, depth)
	n := runtime.Callers(skip+2, fpcs) // +2的目的是跳出GenStack和Callers本身
	j := 0
	for i := n - 1; i >= 0; i-- {
		cs := CallerInfo(fpcs[i] - 1)
		if len(cs) == 0 {
			continue
		}

		if j > 0 {
			buf.WriteString(" --> ")
		}
		buf.WriteString(cs)
		j++
	}

	return buf.String()
}

// PanicFunc panic时可以执行传入的函数
func PanicFunc(panicDo func(panic string)) {
	err := recover()
	if err != nil {
		// 1的目的是跳过Panic本身
		panicInfo := fmt.Sprintf("panic:%v,stack:%s\n", err, CallStack(1, 10))
		panicDo(panicInfo)
	}
}

// PanicDeal 处理recover错误情况
func PanicDeal(panicErr interface{}, panicDo func(panic string)) {
	if panicErr != nil {
		// 1的目的是跳过Panic本身
		panicInfo := fmt.Sprintf("panic:%v,stack:%s\n", panicErr, CallStack(1, 10))
		panicDo(panicInfo)
	}
}

// GetCallerName ...
func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return GetFuncName(runtime.FuncForPC(pc).Name())
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
