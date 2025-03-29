package service

import (
	"aDi/config"
	"aDi/log"
	"os"
	"testing"
)

func init() {
	log.Init(true, "../debug.log")
	config.InitStaticConf()
}

func TestCosDoUpload(t *testing.T) {
	imgBuff, _ := os.ReadFile("E:\\Downloads\\images.jpg")
	url, err := CosDoUpload(imgBuff, "/voice/3801728859/tmp.jpg", nil)
	log.Debug(url, err)
}
