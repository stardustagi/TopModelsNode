package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UsersKey struct {
	Id          int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	UserId      int64  `json:"user_id" xorm:"'user_id' BIGINT(12)"`
	AccessKey   string `json:"access_key" xorm:"'access_key' VARCHAR(128)"`
	SecurityKey string `json:"security_key" xorm:"'security_key' VARCHAR(128)"`
	CreatedAt   int64  `json:"created_at" xorm:"'created_at' comment('创建时间') BIGINT(255)"`
	Deleted     int64  `json:"deleted" xorm:"'deleted' comment('删除时间') BIGINT(1)"`
}

func (o *UsersKey) TableName() string {
	return "users_key"
}

func (o *UsersKey) GetSliceName(slice string) string {
	return fmt.Sprintf("users_key_%s", slice)
}

func (o *UsersKey) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("users_key_%d%02d", t.Year(), t.Month())
}

func (o *UsersKey) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("users_key_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UsersKey) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UsersKey) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UsersKey) PrimaryKey() interface{} {
	return o.Id
}
