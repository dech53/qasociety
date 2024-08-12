package mail

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

// SendEmailCode 发送验证码
func SendEmailCode(code, email, title string) error {
	message := "你的验证码为" + code + ",3分钟内有效"
	host := "smtp.qq.com"
	port := 25
	userName := "2496916936@qq.com"
	//授权码
	passWord := ""
	m := gomail.NewMessage()
	m.SetHeader("From", userName)
	m.SetHeader("From", "dech53"+"<"+userName+">")
	//测试邮箱填自己的就行
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
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
