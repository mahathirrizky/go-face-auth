package helper

import (
	"bytes"
	"fmt"
	"go-face-auth/models"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GenerateInvoicePDF creates a PDF document for an invoice and returns it as a byte slice.
func GenerateInvoicePDF(invoice *models.InvoiceTable) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set Font
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(40, 10, "INVOICE")
	pdf.Ln(20)

	// Company Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Billed To:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, invoice.Company.Name)
	pdf.Ln(5)
	pdf.Cell(40, 10, invoice.Company.Address)
	pdf.Ln(15)

	// Invoice Details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Invoice Number:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, invoice.OrderID)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Payment Date:")
	pdf.SetFont("Arial", "", 12)
	if invoice.PaidAt != nil {
		pdf.Cell(40, 10, invoice.PaidAt.Format("02 January 2006"))
	} else {
		pdf.Cell(40, 10, time.Now().Format("02 January 2006"))
	}
	pdf.Ln(15)

	// Table Header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(130, 10, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 10, "Amount", "1", 0, "R", true, 0, "")
	pdf.Ln(10)

	// Table Body
	pdf.SetFont("Arial", "", 12)
	description := fmt.Sprintf("Subscription for %s Package", invoice.SubscriptionPackage.Name)
	amount := fmt.Sprintf("Rp %.2f", invoice.Amount)
	pdf.CellFormat(130, 10, description, "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 10, amount, "1", 0, "R", false, 0, "")
	pdf.Ln(10)

	// Total
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(130, 10, "Total", "1", 0, "R", false, 0, "")
	pdf.CellFormat(60, 10, amount, "1", 0, "R", false, 0, "")
	pdf.Ln(20)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Thank you for your business!")
	pdf.Ln(5)
	pdf.Cell(0, 10, "This is a computer-generated invoice and does not require a signature.")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
