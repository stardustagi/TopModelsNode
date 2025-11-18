package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ModelsInfo struct {
	Id          int64  `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	ModelId     string `json:"model_id" xorm:"'model_id' not null comment('模型ID') VARCHAR(128)"`
	NodeId      string `json:"node_id" xorm:"'node_id' comment('node编号') VARCHAR(64)"`
	Name        string `json:"name" xorm:"'name' comment('模型名') VARCHAR(128)"`
	ApiVersion  string `json:"api_version" xorm:"'api_version' VARCHAR(24)"`
	DeployName  string `json:"deploy_name" xorm:"'deploy_name' VARCHAR(128)"`
	InputPrice  int    `json:"input_price" xorm:"'input_price' INT(10)"`
	OutputPrice int    `json:"output_price" xorm:"'output_price' INT(10)"`
	CachePrice  int    `json:"cache_price" xorm:"'cache_price' INT(10)"`
	Status      string `json:"status" xorm:"'status' comment('模型状态') VARCHAR(12)"`
	LastUpdate  int64  `json:"last_update" xorm:"'last_update' comment('最后更新时间') BIGINT(20)"`
	IsPrivate   int    `json:"is_private" xorm:"'is_private' comment('是否私有化') TINYINT(255)"`
	OwnerId     int64  `json:"owner_id" xorm:"'owner_id' comment('用户ID') BIGINT(12)"`
	Address     string `json:"address" xorm:"'address' comment('模型地址') VARCHAR(255)"`
}

func (o *ModelsInfo) TableName() string {
	return "models_info"
}

func (o *ModelsInfo) GetSliceName(slice string) string {
	return fmt.Sprintf("models_info_%s", slice)
}

func (o *ModelsInfo) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("models_info_%d%02d", t.Year(), t.Month())
}

func (o *ModelsInfo) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("models_info_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *ModelsInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *ModelsInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *ModelsInfo) PrimaryKey() interface{} {
	return o.Id
}
