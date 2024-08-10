package mail

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendEmailCode(code, email string) error {
	message := "你的验证码为" + code + ",3分钟内有效"
	host := "smtp.qq.com"
	port := 25
	userName := "2496916936@qq.com"
	//授权码
	passWord := ""
	m := gomail.NewMessage()
	m.SetHeader("From", userName)
	m.SetHeader("From", "dech53"+"<"+userName+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "重置密码")
	m.SetBody("text/plain", message)
	d := gomail.NewDialer(
		host,
		port,
		userName,
		passWord,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}
