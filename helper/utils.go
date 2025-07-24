package helper

import (
	"crypto/rand"
	"encoding/base64"
	
	"os"
)

// GenerateRandomString generates a cryptographically secure random string of a given length.
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

// GetFrontendAdminURL retrieves the frontend admin URL from environment variables.
func GetFrontendAdminURL() string {
	url := os.Getenv("FRONTEND_ADMIN_URL")
	if url == "" {
		return "http://admin.localhost:5173" // Default for development
	}
	return url
}
