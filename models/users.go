package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Users struct {
	Id                       int64  `json:"id" xorm:"'id' pk autoincr BIGINT(12)"`
	UserName                 string `json:"user_name" xorm:"'user_name' comment('用户名') VARCHAR(20)"`
	Email                    string `json:"email" xorm:"'email' comment('邮件地址') VARCHAR(32)"`
	Phone                    string `json:"phone" xorm:"'phone' comment('电话号码') VARCHAR(20)"`
	Password                 string `json:"password" xorm:"'password' comment('密码') VARCHAR(128)"`
	RealName                 string `json:"real_name" xorm:"'real_name' comment('真实性名') VARCHAR(20)"`
	IdNumber                 string `json:"id_number" xorm:"'id_number' comment('身份证') VARCHAR(20)"`
	Active                   int    `json:"active" xorm:"'active' comment('是否激活') TINYINT(4)"`
	CreatedAt                int64  `json:"created_at" xorm:"'created_at' comment('创建日期') BIGINT(12)"`
	LastUpdate               int64  `json:"last_update" xorm:"'last_update' comment('最后更新日期') BIGINT(12)"`
	CompanyId                int64  `json:"company_id" xorm:"'company_id' comment('公司Id') BIGINT(12)"`
	Salt                     string `json:"salt" xorm:"'salt' comment('盐') VARCHAR(10)"`
	WalletAddressId          int64  `json:"wallet_address_id" xorm:"'wallet_address_id' comment('钱包地址') BIGINT(12)"`
	SpreadId                 int64  `json:"spread_id" xorm:"'spread_id' comment('推广用户') BIGINT(20)"`
	IsRealnameAuthentication int    `json:"is_realname_authentication" xorm:"'is_realname_authentication' comment('是否实名认证') TINYINT(1)"`
	LastLoginIp              string `json:"last_login_ip" xorm:"'last_login_ip' comment('最后登录IP') VARCHAR(20)"`
	IsBan                    int    `json:"is_ban" xorm:"'is_ban' comment('是否封号') TINYINT(1)"`
	Deleted                  int    `json:"deleted" xorm:"'deleted' comment('是否删除') TINYINT(1)"`
	MailCode                 string `json:"mail_code" xorm:"'mail_code' comment('邮件验证') VARCHAR(128)"`
	PhoneCode                string `json:"phone_code" xorm:"'phone_code' comment('电话验证') VARCHAR(128)"`
	IsAdmin                  int    `json:"is_admin" xorm:"'is_admin' comment('管理员标志') TINYINT(1)"`
	IsPrivate                int    `json:"is_private" xorm:"'is_private' comment('私有化部署用户') TINYINT(1)"`
}

func (o *Users) TableName() string {
	return "users"
}

func (o *Users) GetSliceName(slice string) string {
	return fmt.Sprintf("users_%s", slice)
}

func (o *Users) GetSliceDateMonthTable() string {
	t := time.Now()
	return fmt.Sprintf("users_%d%02d", t.Year(), t.Month())
}

func (o *Users) GetSliceDateDayTable() string {
	t := time.Now()
	return fmt.Sprintf("users_%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func (o *Users) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Users) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *Users) PrimaryKey() interface{} {
	return o.Id
}
