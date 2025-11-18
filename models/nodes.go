package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Nodes struct {
	Id           int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	Ids          string `json:"ids" xorm:"'ids' VARCHAR(24)"`
	NodeUserId   int64  `json:"node_user_id" xorm:"'node_user_id' BIGINT(12)"`
	CreatedAt    int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	LastupdateAt int64  `json:"lastupdate_at" xorm:"'lastupdate_at' BIGINT(12)"`
	Domain       string `json:"domain" xorm:"'domain' VARCHAR(128)"`
}

func (o *Nodes) TableName() string {
	return "nodes"
}

func (o *Nodes) GetSliceName(slice string) string {
	return fmt.Sprintf("nodes_%s", slice)
}

func (o *Nodes) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("nodes_%d%02d", t.Year(), t.Month())
}

func (o *Nodes) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("nodes_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *Nodes) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Nodes) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *Nodes) PrimaryKey() interface{} {
	return o.Id
}
