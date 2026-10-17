package config

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/brothify/internal/models"
	"github.com/google/uuid"
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

	tmpl, err := template.ParseFiles(templatePath("reservation_status.html"))
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

func (e *EmailSender) SendFeedbackReminderEmail(reservation *models.Reservation) error {
	body, err := buildFeedbackReminderHTML(reservation, os.Getenv("FEEDBACK_PAGE_URL"))
	if err != nil {
		return err
	}

	return SendEmail(
		reservation.RESERVATIONPERSONEMAIL,
		"How was your Brothify reservation?",
		body,
	)
}

func templatePath(name string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return name
	}
	return filepath.Join(filepath.Dir(filepath.Dir(filename)), "templates", name)
}

func buildFeedbackReminderHTML(reservation *models.Reservation, pageURL string) (string, error) {
	feedbackURL, err := buildFeedbackURL(pageURL, reservation.ID)
	if err != nil {
		return "", err
	}

	tmpl, err := template.ParseFiles(templatePath("feedback.html"))
	if err != nil {
		return "", fmt.Errorf("failed to load feedback email template: %w", err)
	}
	data := struct {
		UserName        string
		FeedbackURL     string
		ReservationID   string
		ReservationDate string
		ReservationTime string
		CompanyName     string
		Year            string
	}{
		UserName:        reservation.RESERVATIONPERSONNAME,
		FeedbackURL:     feedbackURL,
		ReservationID:   reservation.ID.String(),
		ReservationDate: reservation.RESERVATIONDATE,
		ReservationTime: reservation.RESERVATIONTIME,
		CompanyName:     "Brothify",
		Year:            strconv.Itoa(time.Now().Year()),
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute feedback email template: %w", err)
	}

	return body.String(), nil
}

func buildFeedbackURL(pageURL string, reservationID uuid.UUID) (string, error) {
	pageURL = strings.TrimSpace(pageURL)
	if pageURL == "" {
		return "", fmt.Errorf("FEEDBACK_PAGE_URL environment variable is not set")
	}
	parsedURL, err := url.Parse(pageURL)
	if err != nil || parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return "", fmt.Errorf("FEEDBACK_PAGE_URL must be an absolute HTTP or HTTPS URL")
	}
	query := parsedURL.Query()
	query.Set("reservation_id", reservationID.String())
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
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
