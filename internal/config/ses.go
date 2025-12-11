package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var SES *ses.Client

func InitSES() {
	SES = ses.NewFromConfig(AWSConfig)
}

func SendEmail(to string, subject string, body string) error {

	input := &ses.SendEmailInput{
		Source: aws.String(os.Getenv("AWS_SES_SENDER")),

		Destination: &types.Destination{
			ToAddresses: []string{to},
		},

		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(body),
				},
			},
		},
	}

	_, err := SES.SendEmail(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("SES send failed: %v", err)
	}

	return nil
}

func SendEmailWithInvoice(to string, htmlContent string) error {
	log.Println("Sending email to:", to)
	subject := "Your Invoice from Brothify"
	return SendEmail(to, subject, htmlContent)
}
