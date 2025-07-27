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

// PasswordResetService defines the interface for password reset business logic.
type PasswordResetService interface {
	ForgotPassword(email, userType string) error
	ResetPassword(token, newPassword string) error
}

// passwordResetService is the concrete implementation of PasswordResetService.
type passwordResetService struct {
	adminCompanyRepo    repository.AdminCompanyRepository
	employeeRepo        repository.EmployeeRepository
	passwordResetRepo   repository.PasswordResetRepository
}

// NewPasswordResetService creates a new instance of PasswordResetService.
func NewPasswordResetService(adminCompanyRepo repository.AdminCompanyRepository, employeeRepo repository.EmployeeRepository, passwordResetRepo repository.PasswordResetRepository) PasswordResetService {
	return &passwordResetService{
		adminCompanyRepo:    adminCompanyRepo,
		employeeRepo:        employeeRepo,
		passwordResetRepo:   passwordResetRepo,
	}
}

func (s *passwordResetService) ForgotPassword(email, userType string) error {
	var userID int
	var userEmail string
	var userName string
	var tokenType string
	var resetURL string

	if userType == "admin" {
		admin, err := s.adminCompanyRepo.GetAdminCompanyByEmail(email)
		if err != nil || admin == nil {
			return nil // Don't reveal user existence
		}
		userID = admin.ID
		userEmail = admin.Email
		userName = admin.Email
		tokenType = "admin_password_reset"
		resetURL = fmt.Sprintf("%s/reset-password?token=", helper.GetFrontendAdminBaseURL())
	} else if userType == "employee" {
		employee, err := s.employeeRepo.GetEmployeeByEmail(email)
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

	if err := s.passwordResetRepo.CreatePasswordResetToken(passwordResetToken); err != nil {
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

func (s *passwordResetService) ResetPassword(token, newPassword string) error {
	passwordResetToken, err := s.passwordResetRepo.GetPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}
	if passwordResetToken == nil || passwordResetToken.Used || passwordResetToken.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("invalid or expired password reset token")
	}

	if err := s.passwordResetRepo.MarkPasswordResetTokenAsUsed(passwordResetToken); err != nil {
		// Log error but continue
	}

	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	switch passwordResetToken.TokenType {
	case "admin_password_reset":
		admin, err := s.adminCompanyRepo.GetAdminCompanyByID(passwordResetToken.UserID)
		if err != nil || admin == nil {
			return fmt.Errorf("user not found")
		}
		admin.Password = string(hashedPassword)
		if err := s.adminCompanyRepo.UpdateAdminCompany(admin); err != nil {
			return fmt.Errorf("failed to update admin password: %w", err)
		}
	case "employee_password_reset", "employee_initial_password":
		employee, err := s.employeeRepo.GetEmployeeByID(passwordResetToken.UserID)
		if err != nil || employee == nil {
			return fmt.Errorf("user not found")
		}
		employee.Password = string(hashedPassword)
		if err := s.employeeRepo.UpdateEmployee(employee); err != nil {
			return fmt.Errorf("failed to update employee password: %w", err)
		}
		if err := s.employeeRepo.SetEmployeePasswordSet(uint(employee.ID), true); err != nil {
			// Log error but continue
		}
	default:
		return fmt.Errorf("invalid token type")
	}

	return nil
}