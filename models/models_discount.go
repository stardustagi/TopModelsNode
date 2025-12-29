package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ModelsDiscount struct {
	Id              int64  `json:"id" xorm:"'id' pk BIGINT(12)"`
	ModelId         int64  `json:"model_id" xorm:"'model_id' BIGINT(12)"`
	ModelProviderId int64  `json:"model_provider_id" xorm:"'model_provider_id' BIGINT(12)"`
	Type            string `json:"type" xorm:"'type' comment('percent,money,time') VARCHAR(12)"`
	Value           int64  `json:"value" xorm:"'value' BIGINT(12)"`
	CreatedAt       int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	LastUpdate      int64  `json:"last_update" xorm:"'last_update' BIGINT(12)"`
}

func (o *ModelsDiscount) TableName() string {
	return "models_discount"
}

func (o *ModelsDiscount) GetSliceName(slice string, num uint32) string {
	var hash uint32
	for _, c := range slice {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % num
	return fmt.Sprintf("models_discount_%d", shardIndex)
}

func (o *ModelsDiscount) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("models_discount_%d%02d", t.Year(), t.Month())
}

func (o *ModelsDiscount) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("models_discount_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *ModelsDiscount) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *ModelsDiscount) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *ModelsDiscount) PrimaryKey() interface{} {
	return o.Id
}
