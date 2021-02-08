package sender

import (
	"errors"
	"fmt"
	"regexp"
)

var ErrSender = errors.New("code can be send by either SMS or e-mail")

type Sender struct {
	reSMS  *regexp.Regexp
	reMail *regexp.Regexp
}

func (sdr Sender) Send(to, code string) error {
	if sdr.IsMail(to) {
		return sdr.SendMail(to, code)
	}
	// Unavailable yet!
	// if sdr.IsSMS(to) {
	// 	return sdr.sendSMS(to, code)
	// }
	return ErrSender
}

func (sdr *Sender) SendMail(to, code string) error {
	sa := SMTP{
		From:     "1366723936@qq.com",
		To:       []string{to},
		AuthCode: "wdkidhwkoyuejaia",
		Host:     "smtp.qq.com",
		Port:     587,
	}
	return sa.Send(
		fmt.Sprintf("【绿椰子】登录验证码：%s", code),
		fmt.Sprintf("您的登录验证码为：%s，请勿泄漏给他人，该验证码10分钟内有效。", code),
	)
}

func (sdr *Sender) SendSMS(to, code string) error {
	return nil
}

func (sdr *Sender) IsSMS(to string) bool {
	if sdr.reSMS == nil {
		sdr.reSMS = regexp.MustCompile(`^(13[0-9]|14[57]|15[012356789]|17[0678]|18[0-9])[0-9]{8}$`)
	}
	return sdr.reSMS.MatchString(to)
}

func (sdr *Sender) IsMail(to string) bool {
	if sdr.reMail == nil {
		sdr.reMail = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)
	}
	return sdr.reMail.MatchString(to)
}
