package constants

var (
	NodeJwtExpire       string = "72h"  // 72小时
	NodeUserTokenExpire string = "720h" // 节点用户Token过期时间
	NodeUserJwtExpire   string = "720h" // 节点用户Jwt过期时间
	RedisPrefix         string
	ModelsKeyPrefix     string
)
