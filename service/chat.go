package service

import (
	"aDi/config"
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"bufio"
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"time"
)

// GetChatAiReq 构建聊天ai请求入参
func GetChatAiReq(content string, historyList []*model.ChatMessage) (reqInfo *model.ALiAiReq) {
	message := make([]*model.AiMessage, 0)
	message = append(message, &model.AiMessage{
		Role:    "system",
		Content: "", // 根据用户的数字人信息构建系统提示词
	})
	for i := len(historyList) - 1; i >= 0; i-- {
		message = append(message, &model.AiMessage{
			Role:    "user",
			Content: historyList[i].UMessage,
		})
		message = append(message, &model.AiMessage{
			Role:    "assistant",
			Content: historyList[i].DMessage,
		})
	}
	message = append(message, &model.AiMessage{
		Role:    "user",
		Content: content, // 用户问题构建
	})
	return &model.ALiAiReq{
		Model:    config.GetAiTextAiModel(),
		Messages: message,
		Parameters: &model.AiParameters{
			ALiReqParameters: model.ALiReqParameters{
				ResultFormat: "message",
				Temperature:  0.1,
			},
			PresencePenalty: -0.2,
		},
	}
}

// AiJsonContentGenerate 使用qw进行json文案信息生成 -- 流式生成
func AiJsonContentGenerate(reqInfo *model.ALiAiReq, keyId int64, asyncInfo *model.AiAsyncResult) (message string, aiTokenUsage *model.ALiOpenStreamUsage, errCode model.ErrCode, errMsg model.ErrMsg) {
	errCode = model.ErrCodeSuccess
	reqBytes := util.MarshalWithoutErr(reqInfo)
	// 生成request请求
	request, err := http.NewRequest(http.MethodPost, config.GetAiApiUrl(), bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Errorf("create new http request fail,err:%s", err.Error())
		return message, aiTokenUsage, model.ECHttpDo, errMsg
	}
	// 添加header
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-DashScope-SSE", "enable") // 设置开启流式返回
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.GetAiSecret()))
	response, err := (&http.Client{Timeout: time.Second * time.Duration(30)}).Do(request) // 默认30秒超时
	if err != nil {
		log.Errorf("do http request fail,err:%s", err.Error())
		return message, aiTokenUsage, model.ECHttpDo, errMsg
	}
	defer response.Body.Close()
	// 申明scan
	sc := bufio.NewScanner(response.Body)
	for {
		if !sc.Scan() {
			break
		}
		tmpS := sc.Text()
		if tmpS == "" {
			continue
		}
		// 获取文案信息
		contentTmp, tokenUsage, hasFinish := ALiStreamContentGet(tmpS)
		message += contentTmp
		if tokenUsage != nil && tokenUsage.TotalTokens > 0 {
			aiTokenUsage.TotalTokens = tokenUsage.TotalTokens
			aiTokenUsage.OutputTokens = tokenUsage.OutputTokens
			aiTokenUsage.InputTokens = tokenUsage.InputTokens
		}
		if hasFinish {
			break
		}
		if contentTmp != "" {
			// 每次有新数据就往redis中进行插入
			asyncInfo.AiResult = message
			if keyId > 0 {
				_ = dao.UpdateAsyncInfo(keyId, util.MarshalToStringWithOutErr(asyncInfo))
			}
		}
	}
	if len(message) > 0 {
		return message, aiTokenUsage, errCode, errMsg
	}
	return message, aiTokenUsage, model.ECS2S, "聊天信息生成为空"
}

/*
中间数据 data:{"output":{"choices":[{"message":{"content":"✨","role":"assistant"},"finish_reason":"null"}]},"usage":{"total_tokens":192,"input_tokens":191,"output_tokens":1},"request_id":"72500730-0b9c-97d3-bbcc-44206f56826f"}
结束数据 data:{"output":{"choices":[{"message":{"content":"\n日期：______","role":"assistant"},"finish_reason":"stop"}]},"usage":{"total_tokens":460,"input_tokens":191,"output_tokens":269},"request_id":"dc477ba7-cad3-97b3-9ebe-c2ad814bf2cf"}
*/
// ALiStreamContentGet 阿里模型流式文案信息获取
func ALiStreamContentGet(stream string) (content string, usage *model.ALiOpenStreamUsage, hasFinish bool) {
	if !strings.Contains(stream, "data") {
		// 不包含data的信息直接跳过
		return content, usage, hasFinish
	}
	stream = strings.TrimSpace(strings.TrimPrefix(stream, "data:"))
	if stream == "" {
		return content, usage, hasFinish
	}

	info := &model.ALiOpenStream{}
	err := jsoniter.UnmarshalFromString(stream, info)
	if err != nil {
		log.Errorf("parse fail,info:")
		return content, usage, hasFinish
	}
	if info.GetFinish() {
		log.Debugf("request id:%s", info.GetRequestId())
	}
	return info.GetContent(), info.GetUsage(), info.GetFinish()
}
