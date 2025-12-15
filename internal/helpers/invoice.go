package helpers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/brothify/internal/models"
)

func BuildInvoiceHTML(res *models.Reservation) (string, error) {
	tmpl, err := template.ParseFiles("internal/templates/invoice.html")
	if err != nil {
		log.Println("Error parsing invoice template:", err)
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, res)
	if err != nil {
		log.Println("Error executing invoice template:", err)
		return "", err
	}

	return buf.String(), nil
}

func BuildEmailReservationHTML(res *models.Reservation) (string, error) {
	tmpl, err := template.ParseFiles("templates/reservationconfirmemail.html")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, res)
	if err != nil {
		return "", err
	}	
	return buf.String(), nil
}
	