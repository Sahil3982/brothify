package helpers

import (
	"bytes"
	"html/template"

	"github.com/brothify/internal/models"
)

func BuildInvoiceHTML(res *models.Reservation) (string, error) {
	tmpl, err := template.ParseFiles("templates/invoice.html")
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
