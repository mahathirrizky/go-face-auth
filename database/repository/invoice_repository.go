package repository

import (
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

// CreateInvoice inserts a new Invoice record into the database.
func (r *invoiceRepository) CreateInvoice(invoice *models.InvoiceTable) error {
	result := r.db.Create(invoice)
	if result.Error != nil {
		log.Printf("Error creating invoice: %v", result.Error)
		return result.Error
	}
	log.Printf("Invoice created with ID: %d, OrderID: %s", invoice.ID, invoice.OrderID)
	return nil
}


// GetInvoiceByOrderID retrieves an invoice by its OrderID, preloading related data.
func (r *invoiceRepository) GetInvoiceByOrderID(orderID string) (*models.InvoiceTable, error) {
	var invoice models.InvoiceTable
	err := r.db.Preload("Company").Preload("SubscriptionPackage").Where("order_id = ?", orderID).First(&invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found, let handler decide response
		}
		return nil, err
	}
	return &invoice, nil
}

// GetInvoicesByCompanyID retrieves all invoices for a specific company, ordered by creation date.
func (r *invoiceRepository) GetInvoicesByCompanyID(companyID uint) ([]models.InvoiceTable, error) {
	var invoices []models.InvoiceTable
	err := r.db.Preload("SubscriptionPackage").
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

// UpdateInvoice updates an existing invoice record in the database.
func (r *invoiceRepository) UpdateInvoice(invoice *models.InvoiceTable) error {
	return r.db.Save(invoice).Error
}