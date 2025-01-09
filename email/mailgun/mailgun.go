package emailsMailgun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	emailsType "github.com/huangchunlong818/go-email/email/type"
	"github.com/mailgun/mailgun-go/v4"
	"strconv"
	"time"
)

// Mailgun发送邮件
type Mailgun struct {
	mg     *mailgun.MailgunImpl
	config *emailsType.Config
}

var newMailgun *Mailgun

func GetNewEmailSend(config emailsType.Config) *Mailgun {
	if newMailgun == nil {
		newMailgun = &Mailgun{
			config: &config,
			mg:     mailgun.NewMailgun(config.Domain, config.Secret),
		}
	}
	return newMailgun
}

// 使用mailgun 发送 支持批量发送，也支持单个发送
func (m *Mailgun) Send(params emailsType.SendMailParams) (string, error) {
	var id string
	strParams, err := json.Marshal(params)
	if err != nil {
		return id, errors.New("mailgun发送邮件-JSON格式化参数错误")
	}
	jsonParams := string(strParams)
	//基础检查
	if len(params.Email) < 1 || params.From == "" || params.Subject == "" || (params.Text == "" && params.Html == "") {
		return id, errors.New("mailgun发送邮件-基础参数检查错误：" + jsonParams)
	}
	if len(params.Email) > emailsType.MAILGUN_SEND_MAX {
		return id, errors.New("mailgun发送邮件-超过最大发送人数：" + strconv.Itoa(emailsType.MAILGUN_SEND_MAX) + jsonParams)
	}

	var nowSend *mailgun.MailgunImpl
	//检查是否切换发送域
	if params.Domain != "" && params.Secret != "" {
		nowSend = mailgun.NewMailgun(params.Domain, params.Secret)
	} else {
		nowSend = m.mg
	}

	// 创建邮件消息
	var toName []string
	for keys, email := range params.Email {
		var recipientName string
		if keys < len(params.ToName) && params.ToName[keys] != "" {
			// 如果存在对应的收件人名称且不为空，使用该名称
			recipientName = fmt.Sprintf("%s <%s>", params.ToName[keys], email)
		} else {
			// 否则，仅使用邮箱
			recipientName = email
		}
		toName = append(toName, recipientName) // 将收件人字符串添加到列表中
	}
	message := nowSend.NewMessage(
		fmt.Sprintf("%s<%s>", params.FromName, params.From), // 发件人信息
		params.Subject,                                      // 邮件主题
		"",                                                  // 邮件正文
		toName...,                                           // 收件人
	)

	if params.Html != "" {
		// 设置邮件消息的 "Content-Type" 为 "text/html" 和设置HTML 发送内容
		message.SetHtml(params.Html)
	} else {
		// 创建邮件消息
		message.SetHtml(params.Text)
	}

	//设置邮件标签
	var (
		tags    []string
		tagsNum int
	)
	if len(params.Tags) > 0 {
		for _, tagValue := range params.Tags {
			if tagValue != "" {
				tags = append(tags, tagValue)
				tagsNum++
			}
		}
		if tagsNum > 3 {
			return id, errors.New("mailgun发送邮件-tag数量最多不能超过3个，参数：" + jsonParams)
		}
		if tagsNum > 0 {
			err = message.AddTag(tags...)
			if err != nil {
				return id, errors.New("mailgun发送邮件-添加TAG失败，参数：" + jsonParams + "错误信息：" + err.Error())
			}
		}
	}

	// 设置邮件标头
	if len(params.Header) > 0 {
		for k, v := range params.Header {
			message.AddHeader(k, v)
		}
	}

	//设置回复邮箱
	if params.ReplyEmail != "" {
		message.SetReplyTo(params.ReplyEmail)
	}

	//执行发送
	var timeOut time.Duration
	if params.Time > 0 {
		timeOut = params.Time
	} else {
		timeOut = time.Second * 10 //默认10秒超时
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	_, id, err = nowSend.Send(ctx, message)
	if err != nil {
		return id, errors.New("mailgun发送邮件-发送失败: " + err.Error() + ", 参数: " + jsonParams)
	}

	return id, nil
}
