package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type NodeUsers struct {
	Id         int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	Email      string `json:"email" xorm:"'email' comment('用户邮件') VARCHAR(128)"`
	Salt       string `json:"salt" xorm:"'salt' VARCHAR(20)"`
	Password   string `json:"password" xorm:"'password' comment('用户密码') VARCHAR(255)"`
	CreatedAt  int64  `json:"created_at" xorm:"'created_at' comment('创建日期') BIGINT(12)"`
	Deleted    int64  `json:"deleted" xorm:"'deleted' comment('删除日期') BIGINT(12)"`
	LastUpdate int64  `json:"last_update" xorm:"'last_update' comment('最后登录日期') BIGINT(12)"`
	IsActive   int    `json:"is_active" xorm:"'is_active' comment('是否激活') TINYINT(1)"`
	CompanyId  int64  `json:"company_id" xorm:"'company_id' comment('公司ID') BIGINT(12)"`
}

func (o *NodeUsers) TableName() string {
	return "node_users"
}

func (o *NodeUsers) GetSliceName(slice string) string {
	return fmt.Sprintf("node_users_%s", slice)
}

func (o *NodeUsers) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("node_users_%d%02d", t.Year(), t.Month())
}

func (o *NodeUsers) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("node_users_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *NodeUsers) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *NodeUsers) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *NodeUsers) PrimaryKey() interface{} {
	return o.Id
}
