package email

import (
	"gopkg.in/gomail.v2"
)

type Emailer interface {
	Send(to string, subject string, content string) error
}

type OutlookClient struct {
	from     string
	password string
}

func NewOutlookClient(from string, password string) OutlookClient {
	return OutlookClient{
		from:     from,
		password: password,
	}
}

func (ec OutlookClient) Send(to string, subject string, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ec.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewDialer("smtp.office365.com", 587, ec.from, ec.password)
	return d.DialAndSend(m)
}
