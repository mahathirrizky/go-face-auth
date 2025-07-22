package handlers

import (
	"fmt"
	"go-face-auth/services"
	"go-face-auth/helper"

	"go-face-auth/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlePaymentConfirmation processes Midtrans payment notifications.
func HandlePaymentConfirmation(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
	var notification helper.MidtransNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		log.Printf("[ERROR] HandlePaymentConfirmation - Failed to bind JSON: %v", err)
		helper.SendError(c, http.StatusBadRequest, "Invalid notification body")
		return
	}

	if err := services.ProcessPaymentConfirmation(notification, hub); err != nil {
			log.Printf("[ERROR] HandlePaymentConfirmation - Error processing payment confirmation for OrderID %s: %v", notification.OrderID, err)
			helper.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		helper.SendSuccess(c, http.StatusOK, "Payment notification processed successfully.", nil)
	}
}

// CreateMidtransTransactionRequest defines the structure for creating a Midtrans transaction.
type CreateMidtransTransactionRequest struct {
	CompanyID             int    `json:"company_id" binding:"required"`
	SubscriptionPackageID int    `json:"subscription_package_id" binding:"required"`
	BillingCycle          string `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
}

// CreateMidtransTransaction handles the creation of a Midtrans Snap transaction.
func CreateMidtransTransaction(c *gin.Context) {
	var req CreateMidtransTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := services.CreateMidtransTransaction(req.CompanyID, req.SubscriptionPackageID, req.BillingCycle)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccess(c, http.StatusOK, "Midtrans transaction created successfully", result)
}

// GetCompanyInvoices handles retrieving all invoices for the authenticated company.
func GetCompanyInvoices(c *gin.Context) {
	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}
	compID, ok := companyID.(float64)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, "Invalid company ID type in token claims.")
		return
	}

	invoices, err := services.GetCompanyInvoices(uint(compID))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve invoices.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Invoices retrieved successfully.", invoices)
}

// DownloadInvoicePDF handles generating and sending an invoice PDF.
func DownloadInvoicePDF(c *gin.Context) {
	orderID := c.Param("order_id")

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}
	compID, _ := companyID.(float64)

	pdfBytes, err := services.DownloadInvoicePDF(orderID, uint(compID))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Invoice-%s.pdf", orderID))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetInvoiceByOrderID handles retrieving invoice details by OrderID (checkout_id).
func GetInvoiceByOrderID(c *gin.Context) {
	orderID := c.Param("order_id") // This will be the checkout_id

	invoice, err := services.GetInvoiceByOrderID(orderID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve invoice.")
		return
	}
	if invoice == nil {
		helper.SendError(c, http.StatusNotFound, "Invoice not found.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Invoice retrieved successfully.", invoice)
}
