package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go/aws"
)

type SesMailer struct {
	client *ses.Client
}

func NewSesMailer(ctx context.Context) (*SesMailer, error) {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	sesClient := ses.NewFromConfig(conf)

	return &SesMailer{client: sesClient}, nil
}

func (m *SesMailer) SendEmail(ctx context.Context, email string, subject string, body string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &types.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(email),
	}

	if _, err := m.client.SendEmail(ctx, input); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
