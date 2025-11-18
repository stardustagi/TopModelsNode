package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserPayLog struct {
	Id          int64  `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UserId      int64  `json:"user_id" xorm:"'user_id' BIGINT(20)"`
	PayAmount   int64  `json:"pay_amount" xorm:"'pay_amount' comment('充值金额，单位分') BIGINT(18)"`
	PayTime     int64  `json:"pay_time" xorm:"'pay_time' comment('充值时间') BIGINT(18)"`
	PayChannel  string `json:"pay_channel" xorm:"'pay_channel' comment('充值渠道') VARCHAR(32)"`
	PayReason   string `json:"pay_reason" xorm:"'pay_reason' comment('支付原因') VARCHAR(255)"`
	AdminUserId int64  `json:"admin_user_id" xorm:"'admin_user_id' comment('充值用户ID 系统充值则ID为0') BIGINT(20)"`
}

func (o *UserPayLog) TableName() string {
	return "user_pay_log"
}

func (o *UserPayLog) GetSliceName(slice string) string {
	return fmt.Sprintf("user_pay_log_%s", slice)
}

func (o *UserPayLog) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_pay_log_%d%02d", t.Year(), t.Month())
}

func (o *UserPayLog) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_pay_log_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserPayLog) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserPayLog) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserPayLog) PrimaryKey() interface{} {
	return o.Id
}
