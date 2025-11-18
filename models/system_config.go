package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type SystemConfig struct {
	Id        int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	Key       string `json:"key" xorm:"'key' comment('配置key') VARCHAR(64)"`
	Value     string `json:"value" xorm:"'value' VARCHAR(255)"`
	CreatedAt int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
}

func (o *SystemConfig) TableName() string {
	return "system_config"
}

func (o *SystemConfig) GetSliceName(slice string) string {
	return fmt.Sprintf("system_config_%s", slice)
}

func (o *SystemConfig) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("system_config_%d%02d", t.Year(), t.Month())
}

func (o *SystemConfig) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("system_config_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *SystemConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *SystemConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *SystemConfig) PrimaryKey() interface{} {
	return o.Id
}
