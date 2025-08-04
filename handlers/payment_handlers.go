package handlers

import (
	"fmt"
	"go-face-auth/services"
	"go-face-auth/helper"

	"go-face-auth/websocket"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// PaymentHandler defines the interface for payment related handlers.
type PaymentHandler interface {
	HandlePaymentConfirmation(hub *websocket.Hub) gin.HandlerFunc
	CreateMidtransTransaction(c *gin.Context)
	GetCompanyInvoices(c *gin.Context)
	DownloadInvoicePDF(c *gin.Context)
	GetInvoiceByOrderID(c *gin.Context)
}

// paymentHandler is the concrete implementation of PaymentHandler.
type paymentHandler struct {
	paymentService services.PaymentService
}

// NewPaymentHandler creates a new instance of PaymentHandler.
func NewPaymentHandler(paymentService services.PaymentService) PaymentHandler {
	return &paymentHandler{
		paymentService: paymentService,
	}
}

// HandlePaymentConfirmation processes Midtrans payment notifications.
func (h *paymentHandler) HandlePaymentConfirmation(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification helper.MidtransNotification
		if err := c.ShouldBindJSON(&notification); err != nil {
			log.Printf("[ERROR] HandlePaymentConfirmation - Failed to bind JSON: %v", err)
			helper.SendError(c, http.StatusBadRequest, "Invalid notification body")
			return
		}

		if err := h.paymentService.ProcessPaymentConfirmation(notification, hub); err != nil {
			// If the error is specifically "invoice not found", it's likely a test notification.
			// Log it for info, but return 200 OK so Midtrans validation passes.
			if strings.Contains(err.Error(), "record not found") {
				log.Printf("[INFO] HandlePaymentConfirmation - Received test notification for OrderID %s. No real invoice found, which is expected for a test.", notification.OrderID)
				helper.SendSuccess(c, http.StatusOK, "Test notification received successfully.", nil)
				return
			}
			
			// For all other errors, return a 500 status code.
			log.Printf("[ERROR] HandlePaymentConfirmation - Error processing payment confirmation for OrderID %s: %v", notification.OrderID, err)
			helper.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		
		helper.SendSuccess(c, http.StatusOK, "Payment notification processed successfully.", nil)
	}
}

// CreateMidtransTransactionRequest defines the structure for creating a Midtrans transaction.
type CreateMidtransTransactionRequest struct {
	CompanyID             int     `json:"company_id" binding:"required"`
	SubscriptionPackageID int     `json:"subscription_package_id"` // Optional if custom_offer_token is provided
	BillingCycle          string  `json:"billing_cycle"`           // Optional if custom_offer_token is provided
	CustomOfferToken      string  `json:"custom_offer_token"`      // Optional: Token for a custom offer
}

// CreateMidtransTransaction handles the creation of a Midtrans Snap transaction.
func (h *paymentHandler) CreateMidtransTransaction(c *gin.Context) {
	var req CreateMidtransTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.paymentService.CreateMidtransTransaction(req.CompanyID, req.SubscriptionPackageID, req.BillingCycle, req.CustomOfferToken)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SendSuccess(c, http.StatusOK, "Midtrans transaction created successfully", result)
}

// GetCompanyInvoices handles retrieving all invoices for the authenticated company.
func (h *paymentHandler) GetCompanyInvoices(c *gin.Context) {
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

	invoices, err := h.paymentService.GetCompanyInvoices(uint(compID))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve invoices.")
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Invoices retrieved successfully.", invoices)
}

// DownloadInvoicePDF handles generating and sending an invoice PDF.
func (h *paymentHandler) DownloadInvoicePDF(c *gin.Context) {
	orderID := c.Param("order_id")

	companyID, exists := c.Get("companyID")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, "Company ID not found in token claims.")
		return
	}
	compID, _ := companyID.(float64)

	pdfBytes, err := h.paymentService.DownloadInvoicePDF(orderID, uint(compID))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Invoice-%s.pdf", orderID))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetInvoiceByOrderID handles retrieving invoice details by OrderID (checkout_id).
func (h *paymentHandler) GetInvoiceByOrderID(c *gin.Context) {
	orderID := c.Param("order_id") // This will be the checkout_id

	invoice, err := h.paymentService.GetInvoiceByOrderID(orderID)
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