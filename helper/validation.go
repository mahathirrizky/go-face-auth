package helper

import (
	"regexp"
)

// emailRegex is a simple regex for email validation.
// This regex is a common compromise for general email validation.
// For stricter validation, consider using a dedicated library or more complex regex.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the provided email string has a valid format.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
