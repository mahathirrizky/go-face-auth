package repository

import "go-face-auth/models"

// InvoiceRepository defines the contract for invoice-related database operations.
type InvoiceRepository interface {
	CreateInvoice(invoice *models.InvoiceTable) error
	GetInvoiceByOrderID(orderID string) (*models.InvoiceTable, error)
	GetInvoicesByCompanyID(companyID uint) ([]models.InvoiceTable, error)
	UpdateInvoice(invoice *models.InvoiceTable) error
}
