package sender

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"sync"
)

const contentType = "Content-Type:text/plain;charset=UTF-8\r\n"

type SMTP struct {
	From     string
	To       []string
	AuthCode string
	Host     string
	Port     int16
}

func (s SMTP) Auth() smtp.Auth {
	return smtp.PlainAuth("", s.From, s.AuthCode, s.Host)
}

func (s SMTP) addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s SMTP) Send(subject, body string) error {
	message := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
		strings.Join(s.To, ","), s.From, subject, contentType, body,
	))
	return smtp.SendMail(s.addr(), s.Auth(), s.From, s.To, message)
}

func (s SMTP) SendHTML(subject, path string, data interface{}) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	buf := Buf()
	if err := t.Execute(buf, &data); err != nil {
		return err
	}
	return s.Send(subject, buf.String())
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func Buf() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func Put(b *bytes.Buffer) {
	b.Reset()
	bufferPool.Put(b)
}
