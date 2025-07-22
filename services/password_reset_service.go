package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func ForgotPassword(email, userType string) error {
	var userID int
	var userEmail string
	var userName string
	var tokenType string
	var resetURL string

	if userType == "admin" {
		admin, err := repository.GetAdminCompanyByEmail(email)
		if err != nil || admin == nil {
			return nil // Don't reveal user existence
		}
		userID = admin.ID
		userEmail = admin.Email
		userName = admin.Email
		tokenType = "admin_password_reset"
		resetURL = fmt.Sprintf("%s/reset-password?token=", helper.GetFrontendAdminBaseURL())
	} else if userType == "employee" {
		employee, err := repository.GetEmployeeByEmail(email)
		if err != nil || employee == nil {
			return nil // Don't reveal user existence
		}
		userID = employee.ID
		userEmail = employee.Email
		userName = employee.Name
		tokenType = "employee_password_reset"
		resetURL = fmt.Sprintf("%s/employee-reset-password?token=", helper.GetFrontendBaseURL())
	} else {
		return fmt.Errorf("invalid user type")
	}

	tokenString := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	passwordResetToken := &models.PasswordResetTokenTable{
		Token:     tokenString,
		UserID:    userID,
		TokenType: tokenType,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := repository.CreatePasswordResetToken(passwordResetToken); err != nil {
		return fmt.Errorf("failed to generate password reset link: %w", err)
	}

	resetURL += tokenString

	go func() {
		if err := helper.SendPasswordResetEmail(userEmail, userName, resetURL); err != nil {
			log.Printf("Error sending password reset email to %s: %v", userEmail, err)
		}
	}()

	return nil
}

func ResetPassword(token, newPassword string) error {
	passwordResetToken, err := repository.GetPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}
	if passwordResetToken == nil || passwordResetToken.Used || passwordResetToken.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("invalid or expired password reset token")
	}

	if err := repository.MarkPasswordResetTokenAsUsed(passwordResetToken); err != nil {
		// Log error but continue
	}

	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	switch passwordResetToken.TokenType {
	case "admin_password_reset":
		admin, err := repository.GetAdminCompanyByID(passwordResetToken.UserID)
		if err != nil || admin == nil {
			return fmt.Errorf("user not found")
		}
		admin.Password = string(hashedPassword)
		if err := repository.UpdateAdminCompany(admin); err != nil {
			return fmt.Errorf("failed to update admin password: %w", err)
		}
	case "employee_password_reset", "employee_initial_password":
		employee, err := repository.GetEmployeeByID(passwordResetToken.UserID)
		if err != nil || employee == nil {
			return fmt.Errorf("user not found")
		}
		employee.Password = string(hashedPassword)
		if err := repository.UpdateEmployee(employee); err != nil {
			return fmt.Errorf("failed to update employee password: %w", err)
		}
		if err := repository.SetEmployeePasswordSet(uint(employee.ID), true); err != nil {
			// Log error but continue
		}
	default:
		return fmt.Errorf("invalid token type")
	}

	return nil
}
