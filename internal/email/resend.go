package email

import "github.com/resend/resend-go/v2"

type ResendClient struct {
	from   string
	client *resend.Client
}

func NewResendClient(from string, apiKey string) ResendClient {
	return ResendClient{
		from:   from,
		client: resend.NewClient(apiKey),
	}
}

func (rc ResendClient) Send(to string, subject string, content string) error {
	params := &resend.SendEmailRequest{
		From:    rc.from,
		To:      []string{to},
		Subject: subject,
		Html:    content,
	}

	_, err := rc.client.Emails.Send(params)
	return err
}
