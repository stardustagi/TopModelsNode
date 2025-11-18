package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type NodeKeys struct {
	Id           int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	NodeId       string `json:"node_id" xorm:"'node_id' VARCHAR(128)"`
	AccessKey    string `json:"access_key" xorm:"'access_key' VARCHAR(64)"`
	SecurityKey  string `json:"security_key" xorm:"'security_key' VARCHAR(64)"`
	LastupdateAt int64  `json:"lastupdate_at" xorm:"'lastupdate_at' comment('创建时间') BIGINT(12)"`
	DeletedAt    int64  `json:"deleted_at" xorm:"'deleted_at' comment('删除时间') BIGINT(12)"`
	NodeUserId   int64  `json:"node_user_id" xorm:"'node_user_id' comment('节点所属用户ID') BIGINT(12)"`
	CompanyId    int64  `json:"company_id" xorm:"'company_id' comment('公司ID') BIGINT(12)"`
}

func (o *NodeKeys) TableName() string {
	return "node_keys"
}

func (o *NodeKeys) GetSliceName(slice string) string {
	return fmt.Sprintf("node_keys_%s", slice)
}

func (o *NodeKeys) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("node_keys_%d%02d", t.Year(), t.Month())
}

func (o *NodeKeys) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("node_keys_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *NodeKeys) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *NodeKeys) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *NodeKeys) PrimaryKey() interface{} {
	return o.Id
}
