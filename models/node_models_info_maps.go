package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type NodeModelsInfoMaps struct {
	Id              int64 `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	NodeId          int64 `json:"node_id" xorm:"'node_id' BIGINT(12)"`
	ModelId         int64 `json:"model_id" xorm:"'model_id' BIGINT(12)"`
	ModelProviderId int64 `json:"model_provider_id" xorm:"'model_provider_id' BIGINT(12)"`
	CreatedAt       int64 `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	UpdatedAt       int64 `json:"updated_at" xorm:"'updated_at' BIGINT(12)"`
}

func (o *NodeModelsInfoMaps) TableName() string {
	return "node_models_info_maps"
}

func (o *NodeModelsInfoMaps) GetSliceName(slice string) string {
	return fmt.Sprintf("node_models_info_maps_%s", slice)
}

func (o *NodeModelsInfoMaps) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("node_models_info_maps_%d%02d", t.Year(), t.Month())
}

func (o *NodeModelsInfoMaps) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("node_models_info_maps_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *NodeModelsInfoMaps) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *NodeModelsInfoMaps) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *NodeModelsInfoMaps) PrimaryKey() interface{} {
	return o.Id
}
