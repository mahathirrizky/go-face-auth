package config

import (
	"log"
	"os"
)

var (
	SMTP_SERVER string
	SMTP_PORT   string
	SMTP_USER   string
	SMTP_PASSWORD string
	SMTP_FROM   string
	FrontendBaseURL string
	FrontendAdminBaseURL string
)

func LoadEmailConfig() {
	SMTP_SERVER = os.Getenv("SMTP_SERVER")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	SMTP_USER = os.Getenv("SMTP_USER")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	SMTP_FROM = SMTP_USER // Reverted to use SMTP_USER

	if SMTP_SERVER == "" || SMTP_PORT == "" || SMTP_USER == "" || SMTP_PASSWORD == "" || SMTP_FROM == "" {
		log.Println("WARNING: One or more SMTP environment variables are not set. Email sending may not work.")
	}

	if FrontendBaseURL == "" {
		log.Println("WARNING: FRONTEND_BASE_URL environment variable is not set. Main frontend links may not work correctly.")
	}

	if FrontendAdminBaseURL == "" {
		log.Println("WARNING: FRONTEND_ADMIN_BASE_URL environment variable is not set. Admin password reset links may not work correctly.")
	}
}
