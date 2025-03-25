package util

import "time"

const (
	DefaultTZ        = 28800 // 默认时区为东八区
	DaySecond        = 86400 // 以秒为单位的一天时间
	HourSecond       = 3600  // 一小时对应的秒
	MinuteSecond     = 60    // 一分钟对应的秒
	DayNumFormat     = "20060102"
	TimestampFormatT = "2006-01-02T15:04:05+08:00"
	TimestampFormat  = "2006-01-02 15:04:05"
)

// ParseTimerFromStr 将string状态日期转化成本地时区的time
func ParseTimerFromStr(dateStr, layout string) (t time.Time, err error) {
	return time.ParseInLocation(layout, dateStr, time.FixedZone("", DefaultTZ))
}
