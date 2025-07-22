package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"time"
)

func SetupInitialPassword(token, password string) error {
	if !helper.IsValidPassword(password) {
		return fmt.Errorf("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number")
	}

	passwordResetToken, err := repository.GetPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("failed to set initial password: %w", err)
	}

	if passwordResetToken == nil || passwordResetToken.Used || passwordResetToken.ExpiresAt.Before(time.Now()) || passwordResetToken.TokenType != "employee_initial_password" {
		return fmt.Errorf("invalid or expired initial password setup token")
	}

	if err := repository.MarkPasswordResetTokenAsUsed(passwordResetToken); err != nil {
		// Continue, but log the error. The token is effectively used.
	}

	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	employee, err := repository.GetEmployeeByID(passwordResetToken.UserID)
	if err != nil || employee == nil {
		return fmt.Errorf("user not found")
	}
	employee.Password = string(hashedPassword)
	if err := repository.UpdateEmployee(employee); err != nil {
		return fmt.Errorf("failed to update employee password: %w", err)
	}

	if err := repository.SetEmployeePasswordSet(uint(employee.ID), true); err != nil {
		// Log error but don't block response, as password is set
	}

	return nil
}
