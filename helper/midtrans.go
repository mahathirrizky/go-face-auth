package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-face-auth/config"
	"io/ioutil"
	"net/http"
	"time"
	"crypto/sha512"
	"encoding/hex"
)

const (
	MidtransSnapSandboxURL    = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	MidtransSnapProductionURL = "https://app.midtrans.com/snap/v1/transactions"
	MidtransNotificationSandboxURL = "https://api.sandbox.midtrans.com/v2/notification/verify"
	MidtransNotificationProductionURL = "https://api.midtrans.com/v2/notification/verify"
)

// SnapCreateTransactionReq represents the request body for Midtrans Snap Create Transaction API.
type SnapCreateTransactionReq struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    CustomerDetails    `json:"customer_details,omitempty"`
	ItemDetails        []ItemDetails      `json:"item_details,omitempty"`
	CreditCard         *CreditCard        `json:"credit_card,omitempty"`
	Callbacks          *Callbacks         `json:"callbacks,omitempty"`
}

// TransactionDetails contains details about the transaction.
type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

// CustomerDetails contains details about the customer.
type CustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`
}

// Address contains address details.
type Address struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"` // ISO 3166-1 alpha-3 country code
}

// ItemDetails contains details about an item in the transaction.
type ItemDetails struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Name     string  `json:"name"`
}

// CreditCard contains credit card specific options.
type CreditCard struct {
	Secure bool `json:"secure"`
}

// Callbacks contains URLs for redirecting after payment.
type Callbacks struct {
	Finish string `json:"finish,omitempty"`
	Error  string `json:"error,omitempty"`
	Pending string `json:"pending,omitempty"`
}

// SnapCreateTransactionRes represents the response body from Midtrans Snap Create Transaction API.
type SnapCreateTransactionRes struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// CreateSnapTransaction sends a request to Midtrans Snap API to create a transaction.
func CreateSnapTransaction(req SnapCreateTransactionReq) (*SnapCreateTransactionRes, error) {
	url := MidtransSnapSandboxURL
	if config.MidtransIsProduction {
		url = MidtransSnapProductionURL
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.SetBasicAuth(config.MidtransServerKey, "") // Server Key with empty password
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	res, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Midtrans API returned non-success status code: %d, body: %s", res.StatusCode, string(body))
	}

	var snapRes SnapCreateTransactionRes
	if err := json.Unmarshal(body, &snapRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &snapRes, nil
}

// MidtransNotification represents the structure of a Midtrans notification payload.
type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	PaymentType       string `json:"payment_type"`
	TransactionID     string `json:"transaction_id"`
	TransactionTime   string `json:"transaction_time"`
	FraudStatus       string `json:"fraud_status"`
	// Add other fields as needed
}

// VerifyMidtransNotificationSignature verifies the signature of a Midtrans notification.
func VerifyMidtransNotificationSignature(notification MidtransNotification) bool {
	stringToHash := notification.OrderID + notification.StatusCode + notification.GrossAmount + config.MidtransServerKey
	hasher := sha512.New()
	hasher.Write([]byte(stringToHash))
	hashed := hasher.Sum(nil)
	generatedSignature := hex.EncodeToString(hashed)

	return generatedSignature == notification.SignatureKey
}