package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type AlarmConfig struct {
	Id         int64  `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	Type       string `json:"type" xorm:"'type' comment('token,billing') VARCHAR(12)"`
	Min        int    `json:"min" xorm:"'min' INT(8)"`
	Max        int    `json:"max" xorm:"'max' INT(8)"`
	CreatedAt  int64  `json:"created_at" xorm:"'created_at' BIGINT(255)"`
	Status     int    `json:"status" xorm:"'status' comment('0,1') TINYINT(1)"`
	LastupDate int64  `json:"lastup_date" xorm:"'lastup_date' BIGINT(12)"`
	UserId     int64  `json:"user_id" xorm:"'user_id' BIGINT(12)"`
}

func (o *AlarmConfig) TableName() string {
	return "alarm_config"
}

func (o *AlarmConfig) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("alarm_config_%d", shardIndex)
}

func (o *AlarmConfig) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("alarm_config_%d%02d", t.Year(), t.Month())
}

func (o *AlarmConfig) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("alarm_config_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *AlarmConfig) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *AlarmConfig) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *AlarmConfig) PrimaryKey() interface{} {
	return o.Id
}
