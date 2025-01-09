package emailsSendgrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	emailsType "github.com/huangchunlong818/go-email/email/type"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"strconv"

	"time"
)

// Sendgrid发送邮件
type Sendgrid struct {
	sd     *sendgrid.Client
	config *emailsType.Config
}

var newSendgrid *Sendgrid

func GetNewEmailSend(config emailsType.Config) *Sendgrid {
	if newSendgrid == nil {
		newSendgrid = &Sendgrid{
			config: &config,
			sd:     sendgrid.NewSendClient(config.Secret),
		}
	}
	return newSendgrid
}

// 使用mailgun 发送 支持批量发送，也支持单个发送
func (s *Sendgrid) Send(params emailsType.SendMailParams) (string, error) {
	var id string
	strParams, err := json.Marshal(params)
	if err != nil {
		return id, errors.New("Sendgrid发送邮件-JSON格式化参数错误")
	}
	jsonParams := string(strParams)
	//基础检查
	if len(params.Email) < 1 || params.From == "" || params.Subject == "" || (params.Text == "" && params.Html == "") {
		return id, errors.New("Sendgrid发送邮件-基础参数检查错误：" + jsonParams)
	}
	if len(params.Email) > emailsType.SENDGRID_SEND_MAX {
		return id, errors.New("mailgun发送邮件-超过最大发送人数：" + strconv.Itoa(emailsType.SENDGRID_SEND_MAX) + jsonParams)
	}

	// 创建一个邮件对象
	from := mail.NewEmail(params.FromName, params.From)

	// 创建 toNames 列表
	var toNames []string
	for i, email := range params.Email {
		var recipientName string
		if i < len(params.ToName) && params.ToName[i] != "" {
			recipientName = fmt.Sprintf("%s <%s>", params.ToName[i], email)
		} else {
			recipientName = email
		}
		toNames = append(toNames, recipientName)
	}

	to := mail.NewEmail(toNames[0], params.Email[0]) //第一个发件人
	m := mail.NewV3MailInit(from, params.Subject, to)

	var content *mail.Content
	if params.Html != "" {
		// 设置邮件消息的 "Content-Type" 为 "text/html" 和设置HTML 发送内容
		content = mail.NewContent("text/html", params.Html)
	} else {
		// 创建邮件消息
		content = mail.NewContent("text/plain", params.Text)
	}
	m.AddContent(content)

	//如果是发送多个
	if len(params.Email) > 1 {
		// 创建一个 Personalization 对象
		personalization := mail.NewPersonalization()
		for k, emails := range params.Email {
			if k > 0 {
				personalization.AddTos(mail.NewEmail(toNames[k], emails))
			}
		}
		m.AddPersonalizations(personalization)
	}

	// 设置邮件标头
	if len(params.Header) > 0 {
		m.Headers = params.Header
	}

	//设置邮件标签
	var (
		tags    []string
		tagsNum int
	)
	if len(params.Tags) > 0 {
		for _, tagValue := range params.Tags {
			if tagValue != "" && len(tagValue) <= 128 {
				tags = append(tags, tagValue)
				tagsNum++
			}
		}
		if tagsNum > 10 {
			return id, errors.New("Sendgrid发送邮件-tag数量最多不能超过10个，长度不高于128，参数：" + jsonParams)
		}
		if tagsNum > 0 {
			// 设置邮件的类别作为标签
			m.Categories = params.Tags
		}
	}

	//设置回复邮箱
	if params.ReplyEmail != "" {
		// 设置回复邮箱
		m.ReplyTo = mail.NewEmail("", params.ReplyEmail)
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
	response, err := s.sd.SendWithContext(ctx, m)
	if err != nil {
		return id, errors.New("Sendgrid发送邮件-发送失败: " + err.Error() + ", 参数: " + jsonParams)
	}

	// 检查响应状态码
	//fmt.Printf("Response Status Code: %d\n", response.StatusCode)

	// 打印请求 ID
	if requestID, ok := response.Headers["X-Message-Id"]; ok {
		if len(requestID) > 0 {
			id = requestID[0]
		}
	}

	return id, nil
}
