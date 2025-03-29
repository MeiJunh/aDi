package util

import (
	"aDi/model"
	"math/rand"
	"strconv"
	"time"
)

// GenerateFundingPoolTradeNo 生成资金池的订单号信息
func GenerateFundingPoolTradeNo(poolType model.FundingPoolType) (tradeNo string) {
	partStr := "DT"
	if poolType == model.FPTRe {
		partStr = "RE"
	}
	tradeNo = GenOutTradeNo(partStr)
	return tradeNo
}

// GenOutTradeNo 生成的微信内部订单号
func GenOutTradeNo(partStr string) (tradeNo string) {
	if partStr == "" {
		partStr = "WX"
	}
	dateStr := getTradeNoDatePart()
	timeStr := getTradeNoNanoTimePart()
	tradeNo = dateStr + partStr + timeStr
	tradeNo += getTradeNoNoncePart(32 - len(tradeNo)) // 保有32的长度
	return tradeNo
}

func getTradeNoDatePart() string {
	dateStr := time.Now().Format("20060102150405")
	return dateStr[2:]
}

func getTradeNoNanoTimePart() string {
	nanosecond := time.Now().Nanosecond()
	nanosecondStr := strconv.Itoa(nanosecond)
	if len(nanosecondStr) >= 5 {
		return nanosecondStr[:5]
	}
	for i := len(nanosecondStr); i < 5; i++ {
		nanosecondStr = "0" + nanosecondStr
	}
	return nanosecondStr
}

func getTradeNoNoncePart(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
