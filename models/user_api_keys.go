package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserApiKeys struct {
	Id         int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	UserId     int64  `json:"user_id" xorm:"'user_id' comment('用户ID') BIGINT(12)"`
	LastUpdate int64  `json:"last_update" xorm:"'last_update' comment('最后更新时间') BIGINT(12)"`
	ApiKeys    string `json:"api_keys" xorm:"'api_keys' comment('api表') TEXT"`
}

func (o *UserApiKeys) TableName() string {
	return "user_api_keys"
}

func (o *UserApiKeys) GetSliceName(slice string) string {
	return fmt.Sprintf("user_api_keys_%s", slice)
}

func (o *UserApiKeys) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_api_keys_%d%02d", t.Year(), t.Month())
}

func (o *UserApiKeys) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_api_keys_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserApiKeys) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserApiKeys) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserApiKeys) PrimaryKey() interface{} {
	return o.Id
}
