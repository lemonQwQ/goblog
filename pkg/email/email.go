package email

import (
	"goblog/pkg/config"
	"goblog/pkg/logger"

	"gopkg.in/gomail.v2"
)

var m *gomail.Message

func init() {
	m = gomail.NewMessage()
}

// Send 发送 code验证码 到邮箱 to 中
func Send(to, code string) error {
	// 配置邮件信息
	from := config.GetString("email.from")
	m.SetAddressHeader("From", from, config.GetString("email.sender"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", config.GetString("email.subject"))
	// m.Embed()
	m.SetBody(config.GetString("email.type"), code)

	// 发生邮件
	d := gomail.NewDialer(config.GetString("email.host"), config.GetInt("email.port"), from, config.GetString("email.pwd"))

	err := d.DialAndSend(m)
	if err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
