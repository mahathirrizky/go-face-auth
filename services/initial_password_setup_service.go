package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"time"
)

// InitialPasswordSetupService defines the interface for initial password setup business logic.
type InitialPasswordSetupService interface {
	SetupInitialPassword(token, password string) error
}

// initialPasswordSetupService is the concrete implementation of InitialPasswordSetupService.
type initialPasswordSetupService struct {
	passwordResetRepo repository.PasswordResetRepository
	employeeRepo      repository.EmployeeRepository
}

// NewInitialPasswordSetupService creates a new instance of InitialPasswordSetupService.
func NewInitialPasswordSetupService(passwordResetRepo repository.PasswordResetRepository, employeeRepo repository.EmployeeRepository) InitialPasswordSetupService {
	return &initialPasswordSetupService{
		passwordResetRepo: passwordResetRepo,
		employeeRepo:      employeeRepo,
	}
}

func (s *initialPasswordSetupService) SetupInitialPassword(token, password string) error {
	if !helper.IsValidPassword(password) {
		return fmt.Errorf("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number")
	}

	passwordResetToken, err := s.passwordResetRepo.GetPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("failed to set initial password: %w", err)
	}

	if passwordResetToken == nil || passwordResetToken.Used || passwordResetToken.ExpiresAt.Before(time.Now()) || passwordResetToken.TokenType != "employee_initial_password" {
		return fmt.Errorf("invalid or expired initial password setup token")
	}

	if err := s.passwordResetRepo.MarkPasswordResetTokenAsUsed(passwordResetToken); err != nil {
		// Continue, but log the error. The token is effectively used.
	}

	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	employee, err := s.employeeRepo.GetEmployeeByID(passwordResetToken.UserID)
	if err != nil || employee == nil {
		return fmt.Errorf("user not found")
	}
	employee.Password = string(hashedPassword)
	if err := s.employeeRepo.UpdateEmployee(employee); err != nil {
		return fmt.Errorf("failed to update employee password: %w", err)
	}

	if err := s.employeeRepo.SetEmployeePasswordSet(uint(employee.ID), true); err != nil {
		// Log error but don't block response, as password is set
	}

	return nil
}