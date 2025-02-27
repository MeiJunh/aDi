package util

import (
	"aDi/log"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var customHttpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// DoHttpRequest 外部http统一调用
func DoHttpRequest(ctx context.Context, api string, method string, headers map[string]string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, api, body)
	if err != nil {
		log.Errorf("create new http request fail,api:%s,err:%s", api, err.Error())
		return nil, err
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	// 进行http调用
	response, err := customHttpClient.Do(request)
	if err != nil {
		log.Errorf("do http request fail,api:%s,err:%s", api, err.Error())
		return nil, err
	}
	defer response.Body.Close()
	// 获取返回内容
	v, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("read info from body fail,api:%s,err:%s", api, err.Error())
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		log.Errorf("http response code not ok, code:%d, api:%s", response.StatusCode, api)
		return nil, fmt.Errorf("http response code not ok, code:%d, body:%s", response.StatusCode, string(v))
	}
	return v, err
}
