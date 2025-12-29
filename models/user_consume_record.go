package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserConsumeRecord struct {
	Id               int64  `json:"id" xorm:"'id' pk autoincr comment('主键，自增') BIGINT(20)"`
	UserId           int64  `json:"user_id" xorm:"'user_id' not null comment('用户ID') index BIGINT(11)"`
	NodeId           int64  `json:"node_id" xorm:"'node_id' BIGINT(11)"`
	DiscountAmount   int64  `json:"discount_amount" xorm:"'discount_amount' default 0 comment('折扣数量') BIGINT(20)"`
	TotalConsumed    int64  `json:"total_consumed" xorm:"'total_consumed' default 0 comment('本次使用的币数量') BIGINT(20)"`
	Caller           string `json:"caller" xorm:"'caller' comment('调用方') index VARCHAR(64)"`
	Model            string `json:"model" xorm:"'model' comment('模型') VARCHAR(64)"`
	ModelId          int64  `json:"model_id" xorm:"'model_id' comment('模型id') BIGINT(11)"`
	ActualProvider   string `json:"actual_provider" xorm:"'actual_provider' comment('服务商') VARCHAR(64)"`
	ActualProviderId string `json:"actual_provider_id" xorm:"'actual_provider_id' comment('服务商id') VARCHAR(64)"`
	ConsumeType      string `json:"consume_type" xorm:"'consume_type' default '' comment('消费类型') VARCHAR(255)"`
	TotalCost        int64  `json:"total_cost" xorm:"'total_cost' default 0 comment('成本') BIGINT(20)"`
	CreatedAt        int64  `json:"created_at" xorm:"'created_at' comment('创建时间') BIGINT(20)"`
	UpdatedAt        int64  `json:"updated_at" xorm:"'updated_at' comment('更新时间') BIGINT(20)"`
}

func (o *UserConsumeRecord) TableName() string {
	return "user_consume_record"
}

func (o *UserConsumeRecord) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("user_consume_record_%d", shardIndex)
}

func (o *UserConsumeRecord) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_consume_record_%d%02d", t.Year(), t.Month())
}

func (o *UserConsumeRecord) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_consume_record_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserConsumeRecord) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserConsumeRecord) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserConsumeRecord) PrimaryKey() interface{} {
	return o.Id
}
