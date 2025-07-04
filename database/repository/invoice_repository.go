package repository

import (
	"go-face-auth/database"
	"go-face-auth/models"
	"log"

	"gorm.io/gorm"
)

// CreateInvoice inserts a new Invoice record into the database.
func CreateInvoice(invoice *models.InvoiceTable) error {
	result := database.DB.Create(invoice)
	if result.Error != nil {
		log.Printf("Error creating invoice: %v", result.Error)
		return result.Error
	}
	log.Printf("Invoice created with ID: %d, OrderID: %s", invoice.ID, invoice.OrderID)
	return nil
}

// GetInvoiceByOrderID retrieves an Invoice record by its OrderID.
func GetInvoiceByOrderID(orderID string) (*models.InvoiceTable, error) {
	var invoice models.InvoiceTable
	result := database.DB.Preload("Company").Preload("SubscriptionPackage").Where("order_id = ?", orderID).First(&invoice)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Invoice not found
		}
		log.Printf("Error getting invoice by OrderID %s: %v", orderID, result.Error)
		return nil, result.Error
	}
	return &invoice, nil
}

// UpdateInvoice updates an existing Invoice record in the database.
func UpdateInvoice(invoice *models.InvoiceTable) error {
	result := database.DB.Save(invoice)
	if result.Error != nil {
		log.Printf("Error updating invoice: %v", result.Error)
		return result.Error
	}
	log.Printf("Invoice updated with ID: %d, OrderID: %s", invoice.ID, invoice.OrderID)
	return nil
}
