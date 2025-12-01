package constants

import (
	"fmt"
	"os"
)

var (
	ApplicationName        = "TopModelsNode"
	ApplicationPrefix      = "node"
	Debug             bool = false
)

var (
	WechatOpenIDKey      string = "wechat:openid"
	UserDefaultPassword  string = "FsNyBfcpdtq1p0063MBu"
	PhoneDefaultPassword string = "6R0b8zfcl6rBb5vkeqVj"
	MailDefaultPassword  string = "ajZxPKxauQaHOtpBJF6W"
	ObtenationIterations int    = 3
	CodeKeyExpire        int    = 30
	Domain               string = "https://example.com"
)

var (
	HeaderId string = "id"
)

var (
	AppName    string = "node"
	AppVersion string = "v1"
)

func Init() {
	if os.Getenv("APP_NAME") != "" {
		AppName = os.Getenv("APP_NAME")
	}
	if os.Getenv("APP_VERSION") != "" {
		AppVersion = os.Getenv("APP_VERSION")
	}
	RedisPrefix = fmt.Sprintf("%s:%s", AppName, AppVersion)
	ModelsKeyPrefix = fmt.Sprintf("%s:llm:models", RedisPrefix)
	NodeKeyPrefix = fmt.Sprintf("%s:node", RedisPrefix)
	NodeUserKeyPrefix = fmt.Sprintf("%s:nodeUser", RedisPrefix)
}

// NodeUserTokenKey 节点用户TokenKey
func NodeUserTokenKey(id int64) string {
	return fmt.Sprintf("nodeUserToken:%d", id)
}

func NodeUserJwtKey(id int64) string {
	return fmt.Sprintf("nodeUserJwt:%d", id)
}

func NodeUserGraphVerifyKey(t string) string {
	return fmt.Sprintf("nodeUserGraphVerifyCode:%s", t)
}

func NodeUserPhoneVerifyKey(phone string) string {
	return fmt.Sprintf("nodeUserPhoneVerifyCode:%s", phone)
}

func NodeUserEmailVerifyKey(email string) string {
	return fmt.Sprintf("nodeUserEmailVerifyCode:%s", email)
}
