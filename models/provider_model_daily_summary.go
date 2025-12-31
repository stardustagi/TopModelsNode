package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProviderModelDailySummary struct {
	Id               int64  `json:"id" xorm:"'id' pk autoincr comment('主键，自增') BIGINT(20)"`
	UserId           int64  `json:"user_id" xorm:"'user_id' not null comment('用户ID') index unique(idx_user_provider_model_type_date) BIGINT(20)"`
	ActualProviderId int    `json:"actual_provider_id" xorm:"'actual_provider_id' not null comment('实际服务商ID') unique(idx_user_provider_model_type_date) INT(11)"`
	ModelId          int    `json:"model_id" xorm:"'model_id' not null comment('模型ID') unique(idx_user_provider_model_type_date) INT(11)"`
	ConsumeType      string `json:"consume_type" xorm:"'consume_type' not null comment('消费类型') unique(idx_user_provider_model_type_date) VARCHAR(32)"`
	Date             string `json:"date" xorm:"'date' not null comment('日期YYYY-MM-DD') index unique(idx_user_provider_model_type_date) VARCHAR(10)"`
	TotalConsumed    int64  `json:"total_consumed" xorm:"'total_consumed' default 0 comment('总消费金额') BIGINT(20)"`
	TotalCost        int64  `json:"total_cost" xorm:"'total_cost' default 0 comment('总成本') BIGINT(20)"`
	UpdatedAt        int64  `json:"updated_at" xorm:"'updated_at' comment('更新时间') BIGINT(20)"`
}

func (o *ProviderModelDailySummary) TableName() string {
	return "provider_model_daily_summary"
}

func (o *ProviderModelDailySummary) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("provider_model_daily_summary_%d", shardIndex)
}

func (o *ProviderModelDailySummary) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("provider_model_daily_summary_%d%02d", t.Year(), t.Month())
}

func (o *ProviderModelDailySummary) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("provider_model_daily_summary_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *ProviderModelDailySummary) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *ProviderModelDailySummary) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *ProviderModelDailySummary) PrimaryKey() interface{} {
	return o.Id
}
