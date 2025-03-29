package service

import (
	"aDi/config"
	"aDi/dao"
	"aDi/log"
	"aDi/model"
	"aDi/util"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/medivhzhan/weapp/v2"
	"sync"
	"time"
)

/*
小程序token获取需要使用公共方法 -- 并且需要支持多进程部署 -- 目前仅仅使用mysql来实现该功能
*/

var WxTokenManager *TokenManager

// TokenManager token管理
type TokenManager struct {
	db          *sqlx.DB
	mutex       sync.Mutex
	tokenCache  string
	expireAt    int64
	refreshLock bool
}

type TokenInfo struct {
	Token      string    `json:"token"`
	UpdateTime time.Time `json:"update_time"`
	IsLocked   bool      `json:"is_locked"`
}

// InitTokenManager 初始化token manager
func InitTokenManager(db *sqlx.DB) {
	WxTokenManager = &TokenManager{
		db: db,
	}
	return
}

// GetWxAccessToken 获取微信token
func GetWxAccessToken() (accessToken string, err error) {
	// 获取微信token
	accessToken, err = WxTokenManager.GetToken()
	if err != nil {
		log.Error("GetWxToken GetWxToken err:", err)
		return accessToken, err
	}
	if accessToken == "" {
		// 如果没有报错，但是token又为空,则等待一段时间再次获取
		time.Sleep(time.Millisecond * 500)
		return WxTokenManager.GetToken()
	}
	return accessToken, nil
}

// GetToken 获取token
func (tm *TokenManager) GetToken() (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// 如果缓存的token未过期，直接返回
	if !tm.isTokenExpired() {
		return tm.tokenCache, nil
	}

	// 检查数据库中的token
	tokenInfo, err := dao.GetTokenFromDB()
	if err != nil {
		log.Errorf("get token fail,err:%s", err.Error())
		return "", err
	}

	// 如果数据库中的token未过期且未被锁定，直接使用
	if !isWxTokenExpired(tokenInfo.ExpireAt) && !tokenInfo.IsLocked {
		tm.tokenCache = tokenInfo.Token
		tm.expireAt = tokenInfo.ExpireAt
		return tokenInfo.Token, nil
	}

	// 尝试获取锁进行刷新
	if dao.WxTokenTryLock() {
		defer dao.WxTokenUnlock()

		// 再次检查token是否已被其他服务器更新
		tokenInfo, err = dao.GetTokenFromDB()
		if err != nil {
			log.Errorf("get token fail,err:%s", err.Error())
			return "", err
		}

		if !isWxTokenExpired(tokenInfo.ExpireAt) {
			tm.tokenCache = tokenInfo.Token
			tm.expireAt = tokenInfo.ExpireAt
			return tokenInfo.Token, nil
		}

		// 调用微信接口获取新token
		var newToken string
		var expireAt int64
		newToken, expireAt, err = RefreshTokenFromWechat()
		if err != nil {
			log.Errorf("refresh token fail,err:%s", err.Error())
			return "", err
		}
		if newToken == "" {
			return "", errors.New("get token is nil")
		}
		// 更新数据库
		err = dao.UpdateWxToken(newToken, expireAt)
		if err != nil {
			log.Errorf("update token fail,err:%s", err.Error())
			return "", err
		}

		tm.tokenCache = newToken
		tm.expireAt = expireAt
		return tm.tokenCache, nil
	}

	// 没有抢到锁 -- 自己等待然后再次重试
	return "", nil
}

// 检查token是否过期
func (tm *TokenManager) isTokenExpired() bool {
	return isWxTokenExpired(tm.expireAt)
}

// isWxTokenExpired 微信token是否过期
func isWxTokenExpired(expireAt int64) bool {
	// 如果到期时间到现在只有不到10分钟了，说明可以更新mysql了
	return expireAt-time.Now().Unix() < util.MinuteSecond*10
}

// RefreshTokenFromWechat 从微信服务器刷新token
func RefreshTokenFromWechat() (token string, expireAt int64, err error) {
	// 实现调用微信接口获取新token的逻辑
	// ...
	// 如果是测试环境需要直接调用生产的服务获取 -- 因为测试和生产是同一个access token
	// 如果是生产的则直接调用微信接口进行token 获取
	tokenRsp, err := weapp.GetAccessToken(config.GetAppId(), config.GetAppSecret())
	if err != nil {
		log.Errorf("get access token fail,err:%s", err.Error())
		return token, expireAt, err
	}
	if tokenRsp.ErrCode != 0 {
		rspStr := util.MarshalToStringWithOutErr(tokenRsp)
		log.Errorf("get access token fail,rsp:%s", rspStr)
		return token, expireAt, fmt.Errorf("wx response is no sucess,rsp:%s", rspStr)
	}
	return tokenRsp.AccessToken, int64(tokenRsp.ExpiresIn), nil
}

// GetPhoneByWxCode 根据code获取手机号
func GetPhoneByWxCode(wxDynamicCode string) (phone string, err error) {
	// 获取wx access token
	accessToken, err := GetWxAccessToken()
	if err != nil {
		log.Errorf("get access token fail,err:%s", err.Error())
		return phone, err
	}
	data, _ := jsoniter.Marshal(map[string]string{
		"code": wxDynamicCode,
	})
	rspBuff, err := util.DoHttpRequest(context.Background(), fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken),
		"POST", map[string]string{"Content-Type": "application/json"}, bytes.NewReader(data))
	if err != nil {
		log.Errorf("GetUserPhoneNumber fail,err:%s", err.Error())
		return phone, err
	}
	rsp := &model.GetUserPhoneNumberRsp{}
	err = jsoniter.Unmarshal(rspBuff, rsp)
	if err != nil {
		log.Errorf("json unmarshal fail,err:%s", err.Error())
		return phone, err
	}
	if rsp.Errcode == 40029 || rsp.Errcode == -1 {
		err = fmt.Errorf("get user phone fail,code:%d,msg:%s", rsp.Errcode, rsp.Errmsg)
		return phone, err
	}

	return rsp.PhoneInfo.PhoneNumber, nil
}
