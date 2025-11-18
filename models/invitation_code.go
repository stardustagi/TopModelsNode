package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type InvitationCode struct {
	Id        int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	Code      string `json:"code" xorm:"'code' VARCHAR(12)"`
	Reason    string `json:"reason" xorm:"'reason' VARCHAR(255)"`
	AdminId   int64  `json:"admin_id" xorm:"'admin_id' BIGINT(12)"`
	CreatedAt int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	IsUse     int    `json:"is_use" xorm:"'is_use' TINYINT(1)"`
	UseTime   int64  `json:"use_time" xorm:"'use_time' BIGINT(12)"`
}

func (o *InvitationCode) TableName() string {
	return "invitation_code"
}

func (o *InvitationCode) GetSliceName(slice string) string {
	return fmt.Sprintf("invitation_code_%s", slice)
}

func (o *InvitationCode) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("invitation_code_%d%02d", t.Year(), t.Month())
}

func (o *InvitationCode) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("invitation_code_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *InvitationCode) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *InvitationCode) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *InvitationCode) PrimaryKey() interface{} {
	return o.Id
}
