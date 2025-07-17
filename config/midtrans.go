package config

import (
	"os"
	"log"
)

var (
	MidtransServerKey string
	MidtransClientKey string
	MidtransIsProduction bool
	FRONTEND_ADMIN_BASE_URL string // New: Base URL for the application, used for Midtrans callbacks
)

func LoadMidtransConfig() {
	MidtransServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	MidtransClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	FRONTEND_ADMIN_BASE_URL = os.Getenv("FRONTEND_ADMIN_BASE_URL") // Load APP_BASE_URL
	
	if os.Getenv("MIDTRANS_IS_PRODUCTION") == "true" {
		MidtransIsProduction = true
	} else {
		MidtransIsProduction = false
	}

	if MidtransServerKey == "" || MidtransClientKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY and MIDTRANS_CLIENT_KEY environment variables must be set.")
	}

	if FRONTEND_ADMIN_BASE_URL == "" {
		log.Fatal("APP_BASE_URL environment variable must be set for Midtrans callbacks.")
	}
}
