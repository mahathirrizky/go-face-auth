package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockPDFGenerator is a mock for helper.PDFGenerator
type MockPDFGenerator struct {
	GenerateInvoicePDFFunc func(invoice *models.InvoiceTable) ([]byte, error)
}

// GenerateInvoicePDF mocks the helper function
func (m *MockPDFGenerator) GenerateInvoicePDF(invoice *models.InvoiceTable) ([]byte, error) {
	if m.GenerateInvoicePDFFunc != nil {
		return m.GenerateInvoicePDFFunc(invoice)
	}
	return nil, nil
}

func TestGetCompanyInvoices(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewPaymentService(mocks.InvoiceRepo, nil, nil, nil, nil, nil)

	invoices := []models.InvoiceTable{{ID: 1}, {ID: 2}}

	t.Run("Success", func(t *testing.T) {
		mocks.InvoiceRepo.GetInvoicesByCompanyIDFunc = func(companyID uint) ([]models.InvoiceTable, error) {
			return invoices, nil
		}

		result, err := service.GetCompanyInvoices(1)

		assert.NoError(t, err)
		assert.Equal(t, invoices, result)
	})
}

func TestDownloadInvoicePDF(t *testing.T) {
	mocks := services.NewMockRepositories()
	mockPDFGenerator := new(MockPDFGenerator)
	service := services.NewPaymentService(mocks.InvoiceRepo, nil, nil, nil, nil, mockPDFGenerator)

	invoice := &models.InvoiceTable{ID: 1, CompanyID: 1, Status: "paid"}
	pdfBytes := []byte("mock pdf content")

	t.Run("Success", func(t *testing.T) {
		mocks.InvoiceRepo.GetInvoiceByOrderIDFunc = func(orderID string) (*models.InvoiceTable, error) {
			return invoice, nil
		}
		mockPDFGenerator.GenerateInvoicePDFFunc = func(inv *models.InvoiceTable) ([]byte, error) {
			return pdfBytes, nil
		}

		result, err := service.DownloadInvoicePDF("order123", 1)

		assert.NoError(t, err)
		assert.Equal(t, pdfBytes, result)
	})

	t.Run("Invoice Not Found", func(t *testing.T) {
		mocks.InvoiceRepo.GetInvoiceByOrderIDFunc = func(orderID string) (*models.InvoiceTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.DownloadInvoicePDF("order123", 1)

		assert.Error(t, err)
		assert.Equal(t, "invoice not found", err.Error())
	})

	t.Run("Unauthorized Company", func(t *testing.T) {
		mocks.InvoiceRepo.GetInvoiceByOrderIDFunc = func(orderID string) (*models.InvoiceTable, error) {
			return invoice, nil
		}

		_, err := service.DownloadInvoicePDF("order123", 2) // Different company ID

		assert.Error(t, err)
		assert.Equal(t, "you are not authorized to download this invoice", err.Error())
	})

	t.Run("Invoice Not Paid", func(t *testing.T) {
		invoice.Status = "pending" // Change status to pending
		mocks.InvoiceRepo.GetInvoiceByOrderIDFunc = func(orderID string) (*models.InvoiceTable, error) {
			return invoice, nil
		}

		_, err := service.DownloadInvoicePDF("order123", 1)

		assert.Error(t, err)
		assert.Equal(t, "invoice is not paid yet", err.Error())
		invoice.Status = "paid" // Reset status
	})
}

func TestGetInvoiceByOrderID(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewPaymentService(mocks.InvoiceRepo, nil, nil, nil, nil, nil)

	invoice := &models.InvoiceTable{ID: 1, OrderID: "order123"}

	t.Run("Success", func(t *testing.T) {
		mocks.InvoiceRepo.GetInvoiceByOrderIDFunc = func(orderID string) (*models.InvoiceTable, error) {
			return invoice, nil
		}

		result, err := service.GetInvoiceByOrderID("order123")

		assert.NoError(t, err)
		assert.Equal(t, invoice, result)
	})

	t.Run("Not Found", func(t *testing.T) {
		mocks.InvoiceRepo.GetInvoiceByOrderIDFunc = func(orderID string) (*models.InvoiceTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.GetInvoiceByOrderID("order123")

		assert.Error(t, err)
	})
}
