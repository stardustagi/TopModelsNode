package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ModelsProvider struct {
	Id          int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	ModelId     string `json:"model_id" xorm:"'model_id' comment('模型ID') VARCHAR(128)"`
	NodeId      string `json:"node_id" xorm:"'node_id' VARCHAR(128)"`
	ProviderId  string `json:"provider_id" xorm:"'provider_id' VARCHAR(128)"`
	Type        string `json:"type" xorm:"'type' VARCHAR(64)"`
	Name        string `json:"name" xorm:"'name' VARCHAR(128)"`
	Endpoint    string `json:"endpoint" xorm:"'endpoint' VARCHAR(128)"`
	ApiType     string `json:"api_type" xorm:"'api_type' VARCHAR(64)"`
	ModelName   string `json:"model_name" xorm:"'model_name' VARCHAR(64)"`
	InputPrice  int    `json:"input_price" xorm:"'input_price' INT(10)"`
	OutputPrice int    `json:"output_price" xorm:"'output_price' INT(10)"`
	CachePrice  int    `json:"cache_price" xorm:"'cache_price' INT(10)"`
	ApiKeys     string `json:"api_keys" xorm:"'api_keys' comment('apikeys列表') TEXT"`
	Deleted     int64  `json:"deleted" xorm:"'deleted' BIGINT(12)"`
	LastUpdate  int64  `json:"last_update" xorm:"'last_update' BIGINT(12)"`
}

func (o *ModelsProvider) TableName() string {
	return "models_provider"
}

func (o *ModelsProvider) GetSliceName(slice string) string {
	return fmt.Sprintf("models_provider_%s", slice)
}

func (o *ModelsProvider) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("models_provider_%d%02d", t.Year(), t.Month())
}

func (o *ModelsProvider) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("models_provider_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *ModelsProvider) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *ModelsProvider) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *ModelsProvider) PrimaryKey() interface{} {
	return o.Id
}
