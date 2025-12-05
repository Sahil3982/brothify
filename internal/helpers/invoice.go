package helpers

import (
	"bytes"
	"context"
	"text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Item struct {
	Name  string
	Qty   int
	Price float64
}

type InvoiceData struct {
	InvoiceNumber string
	Date          string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Items         []Item
	TotalAmount   float64
}

func GenerateInvoicePDF(data InvoiceData) ([]byte, error) {
	// Implementation to generate PDF invoice
	temp, err := template.ParseFiles("templates/invoice.html")
	if err != nil {
		return nil, err
	}

	var htmlContent bytes.Buffer

	err = temp.Execute(&htmlContent, data)
	if err != nil {
		return nil, err
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfBytes []byte

	if err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+htmlContent.String()),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBytes, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	); err != nil {
		return nil, err
	}

	return pdfBytes, nil
}
