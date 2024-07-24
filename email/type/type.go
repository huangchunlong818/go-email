package emailType

import "time"

// 请求参数 目前只有 mailgun 用到
type SendMailParams struct {
	From       string            // 发件人信息-以这个为主，主要是发件人邮箱
	FromName   string            // 发件人信息-辅助，发件人名称
	Subject    string            // 邮件主题
	Text       string            // 邮件正文 跟Html需要设置一个就行， 优先Html
	Email      []string          // 收件人邮箱,至少一个，至多1000个
	Html       string            // 设置邮件消息的 "Content-Type" 为 "text/html" 并且设置body  跟Text需要设置一个就行， 优先本参数
	ReplyEmail string            //回复的邮箱，就是收件人回复的邮件收到的邮箱地址
	Header     map[string]string // 设置头部参数，比如退订等
	Tags       []string          //邮件标签
	Time       time.Duration     //发送超时时间，默认10秒
	Domain     string            //如果要切换发送域，这里带上发送域名，只对mailgun起作用
	Secret     string            //如果要切换发送域，这里带上发送域名密钥，只对mailgun起作用
}

// 配置参数
type Config struct {
	Domain, Secret, Type string //发送域名和密钥以及邮件发送类型，目前只有 mailgun 用到
	//其他发送邮件类型配置项
}

const (
	MAILGUN_SEND_MAX  = 1000 //mailgun单次最大发送给多少人数量限制
	SENDGRID_SEND_MAX = 1000 //sendgrid单次最大发送给多少人数量限制
)
