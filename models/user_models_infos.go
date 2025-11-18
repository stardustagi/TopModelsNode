package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserModelsInfos struct {
	Id         int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	NodeId     string `json:"node_id" xorm:"'node_id' comment('节点ID') VARCHAR(128)"`
	ModelIds   string `json:"model_ids" xorm:"'model_ids' comment('模型IDS') TEXT"`
	UserId     int64  `json:"user_id" xorm:"'user_id' comment('用户ID') BIGINT(128)"`
	LastUpdate int64  `json:"last_update" xorm:"'last_update' comment('最后更新时间') BIGINT(12)"`
}

func (o *UserModelsInfos) TableName() string {
	return "user_models_infos"
}

func (o *UserModelsInfos) GetSliceName(slice string) string {
	return fmt.Sprintf("user_models_infos_%s", slice)
}

func (o *UserModelsInfos) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("user_models_infos_%d%02d", t.Year(), t.Month())
}

func (o *UserModelsInfos) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("user_models_infos_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *UserModelsInfos) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *UserModelsInfos) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *UserModelsInfos) PrimaryKey() interface{} {
	return o.Id
}
