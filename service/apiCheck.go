package service

import (
	"aDi/log"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Secret key used for signing the JWT
var jwtKey = []byte("your_secret_key")

// GenerateToken 获取token
func GenerateToken(userID int64) (tokenStr string, expireAt int64, err error) {
	expireAt = time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expireAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(jwtKey)
	if err != nil {
		log.Errorf("jwt sign err: %s", err.Error())
		return "", expireAt, err
	}
	return tokenStr, expireAt, err
}

// ValidateToken token校验
func ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, jwt.ErrInvalidKey
}

// GetSignApiKey 获取参数签名的api key -- 目前直接使用open id后16为作为app key
func GetSignApiKey(openId string) string {
	if len(openId) <= 16 {
		return openId
	}
	return openId[len(openId)-16:]
}
