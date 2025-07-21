package helper

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// emailRegex is a simple regex for email validation.
// This regex is a common compromise for general email validation.
// For stricter validation, consider using a dedicated library or more complex regex.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the provided email string has a valid format.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// GetValidationError extracts a user-friendly error message from a validation error.
func GetValidationError(err error) string {
	if err == nil {
		return ""
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errorMessages = append(errorMessages, fieldError.Field()+" is required")
			case "email":
				errorMessages = append(errorMessages, fieldError.Field()+" must be a valid email address")
			case "min":
				errorMessages = append(errorMessages, fieldError.Field()+" must be at least "+fieldError.Param()+" characters long")
			case "max":
				errorMessages = append(errorMessages, fieldError.Field()+" must be at most "+fieldError.Param()+" characters long")
			case "oneof":
				errorMessages = append(errorMessages, fieldError.Field()+" must be one of "+fieldError.Param())
			default:
				errorMessages = append(errorMessages, fieldError.Field()+" has an invalid value")
			}
		}
		return strings.Join(errorMessages, ", ")
	}

	return err.Error()
}

// IsValidPassword checks if the provided password meets the complexity requirements.
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}