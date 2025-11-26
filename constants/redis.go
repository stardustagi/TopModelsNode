package constants

import "fmt"

var (
	NodeJwtExpire       string = "72h"  // 72小时
	NodeUserTokenExpire string = "720h" // 节点用户Token过期时间
	NodeUserJwtExpire   string = "720h" // 节点用户Jwt过期时间
	RedisPrefix         string
	ModelsKeyPrefix     string
)

var (
	ModelsNodeExpireTime               string = "30s" // 模型节点过期时间
	ModelsNodeExpireTimeString         string = "10m"
	NodeUserTokenExpireTimeString      string = "24h"
	NodeUserMailVerifyExpireTimeString string = "24h"
)

// NodeAccessKey 节点访问Key
func NodeAccessKey(nodeUserId int64, ak string) string {
	return fmt.Sprintf("%d:%s", nodeUserId, ak)
}

// NodeUserMailVerifyKey 节点用户邮箱验证Key
func NodeUserMailVerifyKey(nodeUserId int64) string {
	return fmt.Sprintf("mail:verify:%d", nodeUserId)
}

// NodeAccessModelsKey 节点访问模型Key
func NodeAccessModelsKey(nodeId int64) string {
	return fmt.Sprintf("node:access:%d", nodeId)
}

// NodeUserAccessTokenKey 节点用户访问Token Key
func NodeUserAccessTokenKey(nodeUserId int64) string {
	return fmt.Sprintf("nodeUser:token:%d", nodeUserId)
}

func ModelsNodeKeepLiveKey(nodeId int64) string {
	return fmt.Sprintf("modelsNode:keepLive:%d", nodeId)
}

func NodeInfoKey(nodeId int64) string {
	return fmt.Sprintf("node:info:%d", nodeId)
}
