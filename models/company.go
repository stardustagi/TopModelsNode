package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Company struct {
	Id                       int64  `json:"id" xorm:"'id' pk BIGINT(12)"`
	CompanyName              string `json:"company_name" xorm:"'company_name' comment('公司名') VARCHAR(128)"`
	CreatedDate              int64  `json:"created_date" xorm:"'created_date' comment('创建日期') BIGINT(12)"`
	LicenseId                string `json:"license_id" xorm:"'license_id' comment('营业执照信息') VARCHAR(128)"`
	IsRealnameAuthentication int    `json:"is_realname_authentication" xorm:"'is_realname_authentication' comment('是否实名认证') TINYINT(1)"`
	Address                  string `json:"address" xorm:"'address' comment('公司地址') VARCHAR(255)"`
	Phone                    string `json:"phone" xorm:"'phone' comment('公司电话') VARCHAR(30)"`
	Mail                     string `json:"mail" xorm:"'mail' comment('公司邮件') VARCHAR(64)"`
}

func (o *Company) TableName() string {
	return "company"
}

func (o *Company) GetSliceName(slice string) string {
	return fmt.Sprintf("company_%s", slice)
}

func (o *Company) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("company_%d%02d", t.Year(), t.Month())
}

func (o *Company) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("company_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *Company) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Company) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *Company) PrimaryKey() interface{} {
	return o.Id
}
