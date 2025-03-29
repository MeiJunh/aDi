package service

import (
	"aDi/config"
	"aDi/log"
	"aDi/model"
	"os"
	"testing"
)

func init() {
	log.Init(true, "../debug.log")
	config.Init()
}

func TestCosDoUpload(t *testing.T) {
	imgBuff, _ := os.ReadFile("E:\\Downloads\\images.jpg")
	url, err := CosDoUpload(imgBuff, "/voice/3801728859/tmp.jpg", nil)
	log.Debug(url, err)
}

func TestChat(t *testing.T) {
	// 根据对应的数字人设定进行参数获取
	aiReq := GetChatAiReq("你是谁", nil)
	// 调用ai进行对话生成
	message, _, errCode, errMsg := AiJsonContentGenerate(aiReq, 0, &model.AiAsyncResult{})
	if errCode != model.ErrCodeSuccess {
		log.Errorf("chat create fail,err code:%d,err msg:%s", errCode, errMsg)
		return
	}
	log.Debug(message)
}
