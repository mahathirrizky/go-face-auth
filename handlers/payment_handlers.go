package handlers

import (
	"fmt"
	"go-face-auth/config"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HandlePaymentConfirmation processes Midtrans payment notifications.
func HandlePaymentConfirmation(c *gin.Context) {
	var notification helper.MidtransNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		log.Printf("[ERROR] HandlePaymentConfirmation - Failed to bind JSON: %v", err)
		helper.SendError(c, http.StatusBadRequest, "Invalid notification body")
		return
	}

	log.Printf("[INFO] HandlePaymentConfirmation - Received notification for OrderID: %s, Status: %s, FraudStatus: %s", notification.OrderID, notification.TransactionStatus, notification.FraudStatus)

	// Verify the signature key
	if !helper.VerifyMidtransNotificationSignature(notification) {
		log.Printf("[WARN] HandlePaymentConfirmation - Invalid signature key for OrderID: %s", notification.OrderID)
		helper.SendError(c, http.StatusUnauthorized, "Invalid signature key")
		return
	}

	log.Printf("[INFO] HandlePaymentConfirmation - Signature verified for OrderID: %s", notification.OrderID)

	// Extract order ID from notification
	orderID := notification.OrderID

	// Retrieve invoice based on order ID
	invoice, err := repository.GetInvoiceByOrderID(orderID)
	if err != nil {
		log.Printf("[ERROR] HandlePaymentConfirmation - Failed to retrieve invoice for OrderID %s: %v", orderID, err)
		helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve invoice for order")
		return
	}
	if invoice == nil {
		log.Printf("[WARN] HandlePaymentConfirmation - Invoice not found for OrderID: %s", orderID)
		helper.SendError(c, http.StatusNotFound, "Invoice not found for this order")
		return
	}

	log.Printf("[INFO] HandlePaymentConfirmation - Found invoice ID %d with status %s for OrderID: %s", invoice.ID, invoice.Status, orderID)

	// Update invoice status and payment gateway transaction ID
	invoice.PaymentGatewayTransactionID = notification.TransactionID

	var httpStatus int
	var responseMessage string

	// Process based on transaction status
	switch notification.TransactionStatus {
	case "capture", "settlement":
		log.Printf("[INFO] HandlePaymentConfirmation - Processing 'settlement' for OrderID: %s", orderID)

		// Check fraud status first
		if notification.FraudStatus != "accept" {
			log.Printf("[WARN] HandlePaymentConfirmation - Fraud status is '%s' for OrderID: %s. Not activating subscription.", notification.FraudStatus, orderID)
			httpStatus = http.StatusBadRequest
			responseMessage = "Payment not accepted due to fraud status."
			break
		}
		log.Printf("[INFO] HandlePaymentConfirmation - Fraud status 'accept' for OrderID: %s", orderID)

		// Update invoice status if not already paid
		if invoice.Status == "paid" {
			log.Printf("[INFO] HandlePaymentConfirmation - Invoice status already 'paid' for OrderID: %s. Skipping update.", orderID)
			httpStatus = http.StatusOK
			responseMessage = "Payment successful and subscription activated"
			break
		}

		log.Printf("[INFO] HandlePaymentConfirmation - Updating invoice status to 'paid' for OrderID: %s", orderID)
		invoice.Status = "paid"
		now := time.Now()
		invoice.PaidAt = &now
		if err := repository.UpdateInvoice(invoice); err != nil {
			log.Printf("[ERROR] HandlePaymentConfirmation - Failed to update invoice status to paid for OrderID %s: %v", orderID, err)
			httpStatus = http.StatusInternalServerError
			responseMessage = "Failed to update invoice status to paid"
			break
		}

		// Retrieve company details, preloading AdminCompaniesTable for email sending
		var company models.CompaniesTable
		if err := database.DB.Preload("AdminCompaniesTable").First(&company, invoice.CompanyID).Error; err != nil {
			log.Printf("[ERROR] HandlePaymentConfirmation - Failed to retrieve company for subscription activation for OrderID %s: %v", orderID, err)
			httpStatus = http.StatusInternalServerError
			responseMessage = "Failed to retrieve company for subscription activation"
			break
		}

		// Check if subscription package details are preloaded
		if invoice.SubscriptionPackage.ID == 0 {
			log.Printf("Error: SubscriptionPackage not preloaded for invoice %s", invoice.OrderID)
			httpStatus = http.StatusInternalServerError
			responseMessage = "Subscription package details not found for invoice"
			break
		}

		// Activate company subscription if not already active
		if company.SubscriptionStatus == "active" {
			log.Printf("[INFO] HandlePaymentConfirmation - Company %d subscription already active for OrderID: %s", company.ID, orderID)
			httpStatus = http.StatusOK
			responseMessage = "Payment successful and subscription activated"
			break
		}

		log.Printf("[INFO] HandlePaymentConfirmation - Activating subscription for company %d, OrderID: %s", company.ID, orderID)
		company.SubscriptionStatus = "active"
		company.SubscriptionStartDate = &now
		
		// Calculate SubscriptionEndDate based on the billing cycle from the invoice
		var endDate time.Time
		if invoice.BillingCycle == "yearly" {
			endDate = now.AddDate(1, 0, 0) // 1 year
		} else {
			endDate = now.AddDate(0, 1, 0) // 1 month
		}
		company.SubscriptionEndDate = &endDate

		if err := database.DB.Save(&company).Error; err != nil {
			log.Printf("[ERROR] HandlePaymentConfirmation - Failed to update company subscription for OrderID %s: %v", orderID, err)
			httpStatus = http.StatusInternalServerError
			responseMessage = "Failed to update company subscription status"
			break
		}
		log.Printf("Company %d subscription activated for package %s until %s", company.ID, invoice.SubscriptionPackage.Name, endDate.Format("2006-01-02"))

		// Send invoice PDF via email in a goroutine
		go func() {
			adminEmail := company.AdminCompaniesTable[0].Email // Assuming the first admin is the primary one
			if adminEmail == "" {
				log.Printf("[WARN] No admin email found for company %d to send invoice.", company.ID)
				return
			}

			// Generate PDF
			invoicePDF, err := helper.GenerateInvoicePDF(invoice)
			if err != nil {
				log.Printf("[ERROR] Failed to generate invoice PDF for OrderID %s: %v", invoice.OrderID, err)
				return
			}

			invoiceFileName := fmt.Sprintf("Invoice-%s.pdf", invoice.OrderID)
			if err := helper.SendInvoiceEmail(adminEmail, company.Name, invoiceFileName, invoicePDF); err != nil {
				log.Printf("[ERROR] Failed to send invoice email for OrderID %s to %s: %v", invoice.OrderID, adminEmail, err)
			}
			log.Printf("[INFO] Invoice PDF email sent for OrderID %s to %s", invoice.OrderID, adminEmail)
		}()

		httpStatus = http.StatusOK
		responseMessage = "Payment successful and subscription activated"

	case "pending":
		log.Printf("[INFO] HandlePaymentConfirmation - Processing 'pending' for OrderID: %s", orderID)
		// Payment is pending, update invoice status
		if invoice.Status != "pending" {
			invoice.Status = "pending"
			if err := repository.UpdateInvoice(invoice); err != nil {
				log.Printf("[ERROR] HandlePaymentConfirmation - Failed to update invoice status to pending for OrderID %s: %v", orderID, err)
				httpStatus = http.StatusInternalServerError
				responseMessage = "Failed to update invoice status to pending"
				break
			}
		}
		httpStatus = http.StatusOK
		responseMessage = "Payment pending"

	case "deny", "expire", "cancel":
		log.Printf("[INFO] HandlePaymentConfirmation - Processing 'deny/expire/cancel' for OrderID: %s", orderID)
		// Payment failed or expired, update invoice status
		if invoice.Status != "failed" && invoice.Status != "expired" && invoice.Status != "cancelled" {
			invoice.Status = notification.TransactionStatus // Use Midtrans status directly
			if err := repository.UpdateInvoice(invoice); err != nil {
				log.Printf("[ERROR] HandlePaymentConfirmation - Failed to update invoice status to failed/expired/cancelled for OrderID %s: %v", orderID, err)
				httpStatus = http.StatusInternalServerError
				responseMessage = "Failed to update invoice status to failed/expired/cancelled"
				break
			}
		}
		httpStatus = http.StatusOK
		responseMessage = "Payment failed or expired"

	default:
		log.Printf("[WARN] HandlePaymentConfirmation - Unknown transaction status '%s' for OrderID: %s", notification.TransactionStatus, orderID)
		httpStatus = http.StatusBadRequest
		responseMessage = "Unknown transaction status"
	}

	helper.SendSuccess(c, httpStatus, responseMessage, nil)
	
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

	// Retrieve company and subscription package details
	var company models.CompaniesTable
	// Preload AdminCompaniesTable to get admin email for sending payment link
	if err := database.DB.Preload("AdminCompaniesTable").First(&company, req.CompanyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.SendError(c, http.StatusNotFound, "Company not found")
		} else {
			helper.SendError(c, http.StatusInternalServerError, "Failed to retrieve company")
		}
		return
	}

	var subPackage models.SubscriptionPackageTable
	if err := database.DB.First(&subPackage, req.SubscriptionPackageID).Error; err != nil {
		helper.SendError(c, http.StatusNotFound, "Subscription package not found")
		return	}

	// Calculate amount based on billing cycle
	var amount float64
	if req.BillingCycle == "yearly" {
		amount = subPackage.PriceYearly
	} else {
		amount = subPackage.PriceMonthly
	}

	// Generate Order ID (UUID)
	orderID := uuid.New().String()

	// Calculate Due Date (e.g., 24 hours from now)
	issuedAt := time.Now()
	dueDate := issuedAt.Add(24 * time.Hour)

	// Create Invoice record with pending status
	invoice := &models.InvoiceTable{
		CompanyID:             company.ID,
		SubscriptionPackageID: subPackage.ID,
		OrderID:               orderID,
		Amount:                amount,
		BillingCycle:          req.BillingCycle,
		Status:                "pending",
		IssuedAt:              issuedAt,
		DueDate:               dueDate,
	}

	if err := repository.CreateInvoice(invoice); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	// Prepare Midtrans Snap transaction request
	snapReq := helper.SnapCreateTransactionReq{
		TransactionDetails: helper.TransactionDetails{
			OrderID:    invoice.OrderID,
			GrossAmount: float64(int64(invoice.Amount)), // Convert to int64 then back to float64 for Midtrans
		},
		CustomerDetails: helper.CustomerDetails{
			FirstName: company.Name,
		},
		ItemDetails: []helper.ItemDetails{
			{
				ID:       fmt.Sprintf("PKG-%d", subPackage.ID),
				Price:    float64(int64(amount)), // Use the calculated amount
				Quantity: 1,
				Name:     subPackage.Name,
			},
		},
		Callbacks: &helper.Callbacks{
			Finish:  fmt.Sprintf("%s/payment/finish?order_id=%s", config.AppBaseURL, invoice.OrderID),
			Error:   fmt.Sprintf("%s/payment/error?order_id=%s", config.AppBaseURL, invoice.OrderID),
			Pending: fmt.Sprintf("%s/payment/pending?order_id=%s", config.AppBaseURL, invoice.OrderID),
		},
	}

	if len(company.AdminCompaniesTable) > 0 {
		snapReq.CustomerDetails.Email = company.AdminCompaniesTable[0].Email
	}

	if company.Address != "" {
		snapReq.CustomerDetails.BillingAddress = &helper.Address{
			FirstName: company.Name, // Or admin's name if available
			Address:   company.Address,
			// You might need to parse the address into city, postal code, etc. if needed
			// For now, we'll put the whole address in the Address field.
			CountryCode: "IDN", // Assuming Indonesia
		}
	}

	// Call Midtrans Snap API
	snapRes, err := helper.CreateSnapTransaction(snapReq)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create Midtrans transaction: %v", err))
		return
	}

	// Update invoice with PaymentURL from Midtrans
	invoice.PaymentURL = snapRes.RedirectURL
	if err := repository.UpdateInvoice(invoice); err != nil {
		helper.SendError(c, http.StatusInternalServerError, "Failed to update invoice with payment URL")
		return
	}

	// Send payment link email to admin
	if len(company.AdminCompaniesTable) > 0 {
		adminEmail := company.AdminCompaniesTable[0].Email // Assuming the first admin is the primary one
		go func() {
			if err := helper.SendPaymentLinkEmail(adminEmail, company.Name, snapRes.RedirectURL); err != nil {
				log.Printf("Failed to send payment link email to %s: %v", adminEmail, err)
			}
		}()
	}

	helper.SendSuccess(c, http.StatusOK, "Midtrans transaction created successfully", gin.H{
		"snap_token": snapRes.Token,
		"redirect_url": snapRes.RedirectURL,
		"order_id": invoice.OrderID,
		"invoice_id": invoice.ID,
		"checkout_id": invoice.OrderID,
	})
}

// GetInvoiceByOrderID handles retrieving invoice details by OrderID (checkout_id).
func GetInvoiceByOrderID(c *gin.Context) {
	orderID := c.Param("order_id") // This will be the checkout_id

	invoice, err := repository.GetInvoiceByOrderID(orderID)
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
