package mailgateway

import (
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	// SMTP 配置
	smtpHost := "smtp.263.net"
	smtFrom := "noreplay@aoyin.hk"
	smtpPort := 465                  // 465 SSL端口，如果用587要改成 StartTLS
	smtpPassword := "Aoying2025!@#$" // 这里填 SMTP 密码
	// 创建消息
	m := gomail.NewMessage()
	m.SetHeader("From", smtFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // 也可以用 text/plain

	// 添加附件（可选）
	// m.Attach("test.pdf")

	// 发送
	d := gomail.NewDialer(smtpHost, smtpPort, smtFrom, smtpPassword)
	d.SSL = true // 如果用465必须启用 SSL；587 端口用 TLS

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
