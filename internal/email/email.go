package email

import "github.com/resend/resend-go/v2"

type EmailClient struct {
	from   string
	client *resend.Client
}

func NewEmailClient(apiKey string, from string) EmailClient {
	return EmailClient{
		from:   from,
		client: resend.NewClient(apiKey),
	}
}

func (ec *EmailClient) Send(to []string, subject string, content string) error {
	params := &resend.SendEmailRequest{
		From:    ec.from,
		To:      to,
		Subject: subject,
		Html:    content,
	}

	_, err := ec.client.Emails.Send(params)
	return err
}
