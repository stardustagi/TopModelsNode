package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProviderConsumeSummary struct {
	Id               int64  `json:"id" xorm:"'id' pk autoincr comment('主键，自增') BIGINT(20)"`
	ActualProviderId int    `json:"actual_provider_id" xorm:"'actual_provider_id' not null comment('实际服务商ID') unique unique(idx_provider_month) INT(11)"`
	Month            string `json:"month" xorm:"'month' unique(idx_provider_month) VARCHAR(7)"`
	TotalConsumed    int64  `json:"total_consumed" xorm:"'total_consumed' default 0 comment('总消费金额') BIGINT(20)"`
	TotalCost        int64  `json:"total_cost" xorm:"'total_cost' default 0 comment('总成本') BIGINT(20)"`
	UpdatedAt        int64  `json:"updated_at" xorm:"'updated_at' comment('更新时间') BIGINT(20)"`
}

func (o *ProviderConsumeSummary) TableName() string {
	return "provider_consume_summary"
}

func (o *ProviderConsumeSummary) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("provider_consume_summary_%d", shardIndex)
}

func (o *ProviderConsumeSummary) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("provider_consume_summary_%d%02d", t.Year(), t.Month())
}

func (o *ProviderConsumeSummary) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("provider_consume_summary_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *ProviderConsumeSummary) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *ProviderConsumeSummary) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *ProviderConsumeSummary) PrimaryKey() interface{} {
	return o.Id
}
