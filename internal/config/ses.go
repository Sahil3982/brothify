package config

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"html/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/brothify/internal/models"
)

var SES *ses.Client

func InitSES() {
	SES = ses.NewFromConfig(AWSConfig)
}

type EmailSender struct{}

type ReservationEmailData struct {
	CustomerName    string
	ReservationID   string
	ReservationDate string
	ReservationTime string
	Guest           int
	Status          string
}

func (e *EmailSender) SendReservationStatusEmail(
	reservation *models.Reservation,
) error {

	tmpl, err := template.ParseFiles(
		"internal/templates/reservation_status.html",
	)
	if err != nil {
		return fmt.Errorf("failed to load email template: %w", err)
	}
	data := ReservationEmailData{
		CustomerName:    reservation.RESERVATIONPERSONNAME,
		ReservationID:   reservation.ID.String(),
		ReservationDate: reservation.RESERVATIONDATE,
		ReservationTime: reservation.RESERVATIONTIME,
		Guest:           reservation.NUMBEROFGUESTS,
		Status:          reservation.STATUS,
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return SendEmail(
		reservation.RESERVATIONPERSONEMAIL,
		"Reservation Status Updated",
		body.String(),
	)
}

func SendEmail(to string, subject string, body string) error {
	if SES == nil {
		return fmt.Errorf("SES client is not initialized")
	}

	sender := strings.TrimSpace(os.Getenv("AWS_SES_SENDER"))
	if sender == "" {
		return fmt.Errorf("AWS_SES_SENDER environment variable is not set")
	}
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("recipient email is empty")
	}
	log.Printf("Sending SES email from %s to %s with subject %q", sender, to, subject)

	input := &ses.SendEmailInput{
		Source: aws.String(sender),

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

	_, err := SES.SendEmail(context.Background(), input)
	if err != nil {
		return fmt.Errorf("SES send failed from %q to %q: %w", sender, to, err)
	}
	log.Printf("SES accepted email for %s", to)

	return nil
}

func SendEmailWithInvoice(to *models.Reservation, htmlContent string) error {
	log.Println("Sending email to:", to)
	subject := "Your Invoice from Brothify"
	return SendEmail(to.RESERVATIONPERSONEMAIL, subject, htmlContent)
}
