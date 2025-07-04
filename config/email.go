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
)

func LoadEmailConfig() {
	SMTP_SERVER = os.Getenv("SMTP_SERVER")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	SMTP_USER = os.Getenv("SMTP_USER")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	SMTP_FROM = SMTP_USER

	if SMTP_SERVER == "" || SMTP_PORT == "" || SMTP_USER == "" || SMTP_PASSWORD == "" || SMTP_FROM == "" {
		log.Println("WARNING: One or more SMTP environment variables are not set. Email sending may not work.")
	}
}
