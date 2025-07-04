package models

import "time"

// InvoiceTable represents an invoice for a company's subscription.
type InvoiceTable struct {
	ID                      int        `gorm:"primaryKey" json:"id"`
	CompanyID               int        `json:"company_id"`
	Company                 CompaniesTable `gorm:"foreignKey:CompanyID" json:"-"`
	SubscriptionPackageID   int        `json:"subscription_package_id"`
	SubscriptionPackage     SubscriptionPackageTable `gorm:"foreignKey:SubscriptionPackageID" json:"-"`
	OrderID                 string     `gorm:"unique;not null" json:"order_id"` // Midtrans Order ID
	Amount                  float64    `gorm:"not null" json:"amount"`
	Status                  string     `gorm:"not null" json:"status"` // e.g., pending, paid, failed, expired, cancelled
	PaymentGatewayTransactionID string `json:"payment_gateway_transaction_id,omitempty"`
	PaymentURL              string     `json:"payment_url,omitempty"`
	IssuedAt                time.Time  `json:"issued_at"`
	DueDate                 time.Time  `json:"due_date"`
	PaidAt                  *time.Time `json:"paid_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}
