package constants

import (
	"fmt"
)

var (
	NodeUserKeyPrefix string = fmt.Sprintf("%s:nodeUser", RedisPrefix)
)
