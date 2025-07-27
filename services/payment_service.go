package services

import (
	"fmt"
	"go-face-auth/config"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/websocket"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentService defines the interface for payment related business logic.
type PaymentService interface {
	ProcessPaymentConfirmation(notification helper.MidtransNotification, hub *websocket.Hub) error
	CreateMidtransTransaction(companyID int, subscriptionPackageID int, billingCycle string, customOfferToken string) (map[string]interface{}, error)
	GetCompanyInvoices(companyID uint) ([]models.InvoiceTable, error)
	DownloadInvoicePDF(orderID string, companyID uint) ([]byte, error)
	GetInvoiceByOrderID(orderID string) (*models.InvoiceTable, error)
}

// paymentService is the concrete implementation of PaymentService.
type paymentService struct {
	invoiceRepo          repository.InvoiceRepository
	companyRepo          repository.CompanyRepository
	subscriptionRepo     repository.SubscriptionPackageRepository
	customOfferRepo      repository.CustomOfferRepository
	adminCompanyRepo     repository.AdminCompanyRepository
	db                   *gorm.DB
}

// NewPaymentService creates a new instance of PaymentService.
func NewPaymentService(invoiceRepo repository.InvoiceRepository, companyRepo repository.CompanyRepository, subscriptionRepo repository.SubscriptionPackageRepository, customOfferRepo repository.CustomOfferRepository, adminCompanyRepo repository.AdminCompanyRepository, db *gorm.DB) PaymentService {
	return &paymentService{
		invoiceRepo:          invoiceRepo,
		companyRepo:          companyRepo,
		subscriptionRepo:     subscriptionRepo,
		customOfferRepo:      customOfferRepo,
		adminCompanyRepo:     adminCompanyRepo,
		db:                   db,
	}
}

func (s *paymentService) ProcessPaymentConfirmation(notification helper.MidtransNotification, hub *websocket.Hub) error {
	log.Printf("[INFO] ProcessPaymentConfirmation - Received notification for OrderID: %s, Status: %s, FraudStatus: %s", notification.OrderID, notification.TransactionStatus, notification.FraudStatus)

	if !helper.VerifyMidtransNotificationSignature(notification) {
		return fmt.Errorf("invalid signature key")
	}

	log.Printf("[INFO] ProcessPaymentConfirmation - Signature verified for OrderID: %s", notification.OrderID)

	invoice, err := s.invoiceRepo.GetInvoiceByOrderID(notification.OrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve invoice for OrderID %s: %w", notification.OrderID, err)
	}
	if invoice == nil {
		return fmt.Errorf("invoice not found for this order")
	}

	log.Printf("[INFO] ProcessPaymentConfirmation - Found invoice ID %d with status %s for OrderID: %s", invoice.ID, invoice.Status, notification.OrderID)

	invoice.PaymentGatewayTransactionID = notification.TransactionID

	switch notification.TransactionStatus {
	case "capture", "settlement":
		log.Printf("[INFO] ProcessPaymentConfirmation - Processing 'settlement' for OrderID: %s", notification.OrderID)

		if notification.FraudStatus != "accept" {
			return fmt.Errorf("payment not accepted due to fraud status: %s", notification.FraudStatus)
		}
		log.Printf("[INFO] ProcessPaymentConfirmation - Fraud status 'accept' for OrderID: %s", notification.OrderID)

		if invoice.Status == "paid" {
			log.Printf("[INFO] ProcessPaymentConfirmation - Invoice status already 'paid' for OrderID: %s. Skipping update.", notification.OrderID)
			return nil
		}

		log.Printf("[INFO] ProcessPaymentConfirmation - Current invoice status for OrderID %s: %s", notification.OrderID, invoice.Status)
		log.Printf("[INFO] ProcessPaymentConfirmation - Attempting to update invoice status to 'paid' for OrderID: %s", notification.OrderID)
		invoice.Status = "paid"
		now := time.Now()
		invoice.PaidAt = &now
		if err := s.invoiceRepo.UpdateInvoice(invoice); err != nil {
			return fmt.Errorf("failed to update invoice status to paid for OrderID %s: %w", notification.OrderID, err)
		}
		log.Printf("[INFO] ProcessPaymentConfirmation - Invoice status updated to 'paid' for OrderID %s. New status: %s", notification.OrderID, invoice.Status)

		company, err := s.companyRepo.GetCompanyByID(invoice.CompanyID)
		if err != nil {
			return fmt.Errorf("failed to retrieve company for subscription activation for OrderID %s: %w", notification.OrderID, err)
		}

		var pkgName string
		if invoice.SubscriptionPackageID != 0 {
			subPackage, err := s.subscriptionRepo.GetSubscriptionPackageByID(invoice.SubscriptionPackageID)
			if err != nil {
				return fmt.Errorf("subscription package not found for invoice %s: %w", invoice.OrderID, err)
			}
			pkgName = subPackage.PackageName
		} else {
			// This is a custom offer, package name should be in the invoice itself if needed for logging
			pkgName = "Custom Package"
		}

		if company.SubscriptionStatus == "active" {
			log.Printf("[INFO] ProcessPaymentConfirmation - Company %d subscription already active for OrderID: %s", company.ID, notification.OrderID)
			return nil
		}

		log.Printf("[INFO] ProcessPaymentConfirmation - Activating subscription for company %d, OrderID: %s", company.ID, notification.OrderID)
		company.SubscriptionStatus = "active"
		company.SubscriptionStartDate = &now

		var endDate time.Time
		if invoice.BillingCycle == "yearly" {
			endDate = now.AddDate(1, 0, 0)
		} else {
			endDate = now.AddDate(0, 1, 0)
		}
		company.SubscriptionEndDate = &endDate

		// If this invoice is for a custom offer, update the company's CustomOfferID
		if invoice.CustomOfferID != nil {
			company.CustomOfferID = invoice.CustomOfferID
		}

		if err := s.companyRepo.UpdateCompany(company); err != nil {
			return fmt.Errorf("failed to update company subscription for OrderID %s: %w", notification.OrderID, err)
		}
		log.Printf("Company %d subscription activated for package %s until %s", company.ID, pkgName, endDate.Format("2006-01-02"))

		go hub.BroadcastSuperAdminDashboardUpdate()

		go func() {
			adminCompany, err := s.adminCompanyRepo.GetAdminCompanyByCompanyID(company.ID)
			if err != nil || adminCompany == nil {
				log.Printf("[WARN] No admin email found for company %d to send invoice.", company.ID)
				return
			}
			adminEmail := adminCompany.Email

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

	case "pending":
		log.Printf("[INFO] ProcessPaymentConfirmation - Processing 'pending' for OrderID: %s", notification.OrderID)
		if invoice.Status != "pending" {
			invoice.Status = "pending"
			if err := s.invoiceRepo.UpdateInvoice(invoice); err != nil {
				return fmt.Errorf("failed to update invoice status to pending for OrderID %s: %w", notification.OrderID, err)
			}
		}

	case "deny", "expire", "cancel":
		log.Printf("[INFO] ProcessPaymentConfirmation - Processing 'deny/expire/cancel' for OrderID: %s", notification.OrderID)
		if invoice.Status != "failed" && invoice.Status != "expired" && invoice.Status != "cancelled" {
			invoice.Status = notification.TransactionStatus
			if err := s.invoiceRepo.UpdateInvoice(invoice); err != nil {
				return fmt.Errorf("failed to update invoice status to failed/expired/cancelled for OrderID %s: %w", notification.OrderID, err)
			}
		}

	default:
		return fmt.Errorf("unknown transaction status: %s", notification.TransactionStatus)
	}

	return nil
}

func (s *paymentService) CreateMidtransTransaction(companyID int, subscriptionPackageID int, billingCycle string, customOfferToken string) (map[string]interface{}, error) {
	company, err := s.companyRepo.GetCompanyByID(companyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company not found")
		} else {
			return nil, fmt.Errorf("failed to retrieve company: %w", err)
		}
	}

	var packageName string
	var amount float64

	var subPackageID uint = 0 // Default to 0, will be set if it's a standard package
	var customOfferID *uint // Nullable, will be set if it's a custom offer

	if customOfferToken != "" {
		offer, err := s.customOfferRepo.GetCustomOfferByToken(customOfferToken)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve custom offer: %w", err)
		}
		if offer.CompanyID != uint(companyID) {
			return nil, fmt.Errorf("custom offer does not belong to this company")
		}
		if offer.Status != "pending" {
			return nil, fmt.Errorf("custom offer is not pending or has already been used")
		}

		packageName = offer.PackageName
		amount = offer.FinalPrice
		billingCycle = offer.BillingCycle // Use billing cycle from offer
		customOfferID = &offer.ID // Set CustomOfferID

		// Mark the custom offer as used immediately to prevent double-use
		if err := s.customOfferRepo.MarkCustomOfferAsUsed(customOfferToken); err != nil {
			return nil, fmt.Errorf("failed to mark custom offer as used: %w", err)
		}

	} else {
		subPackage, err := s.subscriptionRepo.GetSubscriptionPackageByID(subscriptionPackageID)
		if err != nil {
			return nil, fmt.Errorf("subscription package not found")
		}
		subPackageID = uint(subPackage.ID)

		packageName = subPackage.PackageName
		if billingCycle == "yearly" {
			amount = subPackage.PriceYearly
		} else {
			amount = subPackage.PriceMonthly
		}

	}

	orderID := uuid.New().String()
	issuedAt := time.Now()
	dueDate := issuedAt.Add(24 * time.Hour)

	invoice := &models.InvoiceTable{
		CompanyID:             company.ID,
		SubscriptionPackageID: int(subPackageID), // This will be 0 if it's a custom offer
		CustomOfferID:         customOfferID,     // Set CustomOfferID
		OrderID:               orderID,
		Amount:                amount,
		BillingCycle:          billingCycle,
		Status:                "pending",
		IssuedAt:              issuedAt,
		DueDate:               dueDate,
	}

	if err := s.invoiceRepo.CreateInvoice(invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	snapReq := helper.SnapCreateTransactionReq{
		TransactionDetails: helper.TransactionDetails{
			OrderID:    invoice.OrderID,
			GrossAmount: float64(int64(invoice.Amount)),
		},
		CustomerDetails: helper.CustomerDetails{
			FirstName: company.Name,
		},
		ItemDetails: []helper.ItemDetails{
			{
				ID:       fmt.Sprintf("PKG-%d", subPackageID),
				Price:    float64(int64(amount)),
				Quantity: 1,
				Name:     packageName,
			},
		},
		Callbacks: &helper.Callbacks{
			Finish:  fmt.Sprintf("%s/payment/finish?order_id=%s", config.FRONTEND_ADMIN_BASE_URL, invoice.OrderID),
			Error:   fmt.Sprintf("%s/payment/error?order_id=%s", config.FRONTEND_ADMIN_BASE_URL, invoice.OrderID),
			Pending: fmt.Sprintf("%s/payment/pending?order_id=%s", config.FRONTEND_ADMIN_BASE_URL, invoice.OrderID),
		},
	}

	if len(company.AdminCompaniesTable) > 0 {
		adminCompany, err := s.adminCompanyRepo.GetAdminCompanyByCompanyID(company.ID)
		if err == nil && adminCompany != nil {
			snapReq.CustomerDetails.Email = adminCompany.Email
		}
	}

	if company.Address != "" {
		snapReq.CustomerDetails.BillingAddress = &helper.Address{
			FirstName: company.Name,
			Address:   company.Address,
			CountryCode: "IDN",
		}
	}

	snapRes, err := helper.CreateSnapTransaction(snapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	invoice.PaymentURL = snapRes.RedirectURL
	if err := s.invoiceRepo.UpdateInvoice(invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice with payment URL: %w", err)
	}

	if len(company.AdminCompaniesTable) > 0 {
		adminCompany, err := s.adminCompanyRepo.GetAdminCompanyByCompanyID(company.ID)
		if err == nil && adminCompany != nil {
			adminEmail := adminCompany.Email
			go func() {
				if err := helper.SendPaymentLinkEmail(adminEmail, company.Name, snapRes.RedirectURL); err != nil {
					log.Printf("Failed to send payment link email to %s: %v", adminEmail, err)
				}
			}()
		}
	}

	return gin.H{
		"snap_token": snapRes.Token,
		"redirect_url": snapRes.RedirectURL,
		"order_id": invoice.OrderID,
		"invoice_id": invoice.ID,
		"checkout_id": invoice.OrderID,
	}, nil
}

func (s *paymentService) GetCompanyInvoices(companyID uint) ([]models.InvoiceTable, error) {
	return s.invoiceRepo.GetInvoicesByCompanyID(companyID)
}

func (s *paymentService) DownloadInvoicePDF(orderID string, companyID uint) ([]byte, error) {
	invoice, err := s.invoiceRepo.GetInvoiceByOrderID(orderID)
	if err != nil || invoice == nil {
		return nil, fmt.Errorf("invoice not found")
	}

	if invoice.CompanyID != int(companyID) {
		return nil, fmt.Errorf("you are not authorized to download this invoice")
	}

	if invoice.Status != "paid" {
		return nil, fmt.Errorf("invoice is not paid yet")
	}

	pdfBytes, err := helper.GenerateInvoicePDF(invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice PDF: %w", err)
	}

	return pdfBytes, nil
}

func (s *paymentService) GetInvoiceByOrderID(orderID string) (*models.InvoiceTable, error) {
	return s.invoiceRepo.GetInvoiceByOrderID(orderID)
}