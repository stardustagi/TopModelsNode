package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserAgentTokens struct {
	Id              int64  `json:"id" xorm:"'id' pk BIGINT(12)"`
	UserId          int64  `json:"user_id" xorm:"'user_id' comment('用户ID') BIGINT(12)"`
	UserAgentTokens string `json:"user_agent_tokens" xorm:"'user_agent_tokens' comment('用户端TokenList') TEXT"`
}

func (o *UserAgentTokens) TableName() string {
	return "user_agent_tokens"
}

func (o *UserAgentTokens) GetSliceName(slice string) string {
	return fmt.Sprintf("user_agent_tokens_%s", slice)
}

func (o *UserAgentTokens) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_agent_tokens_%d%02d", t.Year(), t.Month())
}

func (o *UserAgentTokens) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_agent_tokens_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserAgentTokens) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserAgentTokens) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserAgentTokens) PrimaryKey() interface{} {
	return o.Id
}
