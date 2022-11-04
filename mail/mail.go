package mail

import (
	"gopkg.in/gomail.v2"
)

type Email struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (e *Email) Send(subject, body string, to, attachs []string) error {
	m := gomail.NewMessage()
	//设置发件人
	m.SetHeader("From", e.Username)
	//设置发送给多个用户
	m.SetHeader("To", to...)
	//设置邮件主题
	m.SetHeader("Subject", subject)
	//设置邮件正文
	m.SetBody("text/html", body)
	for _, attach := range attachs {
		m.Attach(attach)
	}
	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	err := d.DialAndSend(m)
	return err
}
