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
	FrontendBaseURL string // New: Base URL for the frontend, used for password reset links
)

func LoadEmailConfig() {
	SMTP_SERVER = os.Getenv("SMTP_SERVER")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	SMTP_USER = os.Getenv("SMTP_USER")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	SMTP_FROM = SMTP_USER
	FrontendBaseURL = os.Getenv("FRONTEND_BASE_URL") // Load FRONTEND_BASE_URL

	if SMTP_SERVER == "" || SMTP_PORT == "" || SMTP_USER == "" || SMTP_PASSWORD == "" || SMTP_FROM == "" {
		log.Println("WARNING: One or more SMTP environment variables are not set. Email sending may not work.")
	}

	if FrontendBaseURL == "" {
		log.Println("WARNING: FRONTEND_BASE_URL environment variable is not set. Password reset links may not work correctly.")
	}
}
