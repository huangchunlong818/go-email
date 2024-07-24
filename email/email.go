package email

import (
	Mailgun "github.com/huangchunlong818/go-email/email/mailgun"
	Sendgrid "github.com/huangchunlong818/go-email/email/sendgrid"
	emailType "github.com/huangchunlong818/go-email/email/type"
)

type Emails interface {
	Send(params emailType.SendMailParams) (string, error) //发送邮件
}

// 调用方法发送邮件，默认使用mailgun 发送
func SendMailMethod(config emailType.Config) Emails {
	//返回发送实例
	switch config.Type {
	case "mailgun":
		return Mailgun.GetNewEmailSend(config)
	case "sendgrid":
		return Sendgrid.GetNewEmailSend(config)
	default:
		return Mailgun.GetNewEmailSend(config)
	}
}
