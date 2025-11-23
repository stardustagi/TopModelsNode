package mailgateway

import "testing"

func TestMail263(t *testing.T) {
	// 263邮箱发信测试
	err := SendEmail("yejunqiang@isoisp.com", "测试邮件", "邮件验证码是:12341234")
	if err != nil {
		t.Fatal("发送邮件失败:", err)
	}
	t.Log("发送邮件成功")
}
