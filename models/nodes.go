package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Nodes struct {
	Id           int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	Name         string `json:"name" xorm:"'name' VARCHAR(24)"`
	OwnerId      int64  `json:"owner_id" xorm:"'owner_id' comment('节点所有者ID对应nodeUserID') BIGINT(12)"`
	CreatedAt    int64  `json:"created_at" xorm:"'created_at' BIGINT(12)"`
	LastupdateAt int64  `json:"lastupdate_at" xorm:"'lastupdate_at' BIGINT(12)"`
	Domain       string `json:"domain" xorm:"'domain' VARCHAR(128)"`
	AccessKey    string `json:"access_key" xorm:"'access_key' comment('ak') VARCHAR(255)"`
	SecurityKey  string `json:"security_key" xorm:"'security_key' comment('sk') VARCHAR(255)"`
	CompanyId    int64  `json:"company_id" xorm:"'company_id' comment('企业ID') BIGINT(12)"`
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
