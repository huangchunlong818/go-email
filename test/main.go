package main

import (
	"fmt"
	"github.com/huangchunlong818/go-email/email"
	emailType "github.com/huangchunlong818/go-email/email/type"
)

func main() {

	//mailgun发送邮件
	m := email.SendMailMethod(emailType.Config{
		Domain: "a.mailgun.com",
		Secret: "xxxxx",
		Type:   "mailgun",
	})
	send, err := m.Send(emailType.SendMailParams{
		From:       "",
		FromName:   "",
		Subject:    "",
		Text:       "",
		Email:      nil,
		Html:       "",
		ReplyEmail: "",
		Header:     nil,
		Tags:       nil,
		Time:       0,
	})
	fmt.Println(send, err)

	//切换mailgun发送域名
	ss, err := m.Send(emailType.SendMailParams{
		From:       "",
		FromName:   "",
		Subject:    "",
		Text:       "",
		Email:      nil,
		Html:       "",
		ReplyEmail: "",
		Header:     nil,
		Tags:       nil,
		Time:       0,
		Domain:     "new.com",    //新的发送域名
		Secret:     "new Secret", //新的发送域名的密钥
	})
	fmt.Println(ss, err)

	//Sendgrid发送邮件
	s := email.SendMailMethod(emailType.Config{
		Secret: "xxxxx",
		Type:   "sendgrid",
	})
	tt, err := s.Send(emailType.SendMailParams{
		From:       "",
		FromName:   "",
		Subject:    "",
		Text:       "",
		Email:      nil,
		Html:       "",
		ReplyEmail: "",
		Header:     nil,
		Tags:       nil,
		Time:       0,
	})
	fmt.Println(tt, err)
}
