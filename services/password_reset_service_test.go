package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestForgotPassword(t *testing.T) {
	mockAdminCompanyRepo := new(MockAdminCompanyRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockPasswordResetRepo := new(MockPasswordResetRepository)
	service := services.NewPasswordResetService(mockAdminCompanyRepo, mockEmployeeRepo, mockPasswordResetRepo)

	admin := &models.AdminCompaniesTable{ID: 1, Email: "admin@test.com"}
	employee := &models.EmployeesTable{ID: 2, Email: "employee@test.com", Name: "Test Employee"}

	t.Run("Admin Success", func(t *testing.T) {
		mockAdminCompanyRepo.GetAdminCompanyByEmailFunc = func(email string) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mockPasswordResetRepo.CreatePasswordResetTokenFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}

		err := service.ForgotPassword("admin@test.com", "admin")

		assert.NoError(t, err)
	})

	t.Run("Employee Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByEmailFunc = func(email string) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockPasswordResetRepo.CreatePasswordResetTokenFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}

		err := service.ForgotPassword("employee@test.com", "employee")

		assert.NoError(t, err)
	})

	t.Run("Invalid User Type", func(t *testing.T) {
		err := service.ForgotPassword("test@test.com", "invalid")

		assert.Error(t, err)
		assert.Equal(t, "invalid user type", err.Error())
	})
}

func TestResetPassword(t *testing.T) {
	mockAdminCompanyRepo := new(MockAdminCompanyRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockPasswordResetRepo := new(MockPasswordResetRepository)
	service := services.NewPasswordResetService(mockAdminCompanyRepo, mockEmployeeRepo, mockPasswordResetRepo)

	adminToken := &models.PasswordResetTokenTable{Token: "admin-token", UserID: 1, ExpiresAt: time.Now().Add(time.Hour), TokenType: "admin_password_reset"}
	employeeToken := &models.PasswordResetTokenTable{Token: "employee-token", UserID: 2, ExpiresAt: time.Now().Add(time.Hour), TokenType: "employee_password_reset"}

	admin := &models.AdminCompaniesTable{ID: 1, Email: "admin@test.com"}
	employee := &models.EmployeesTable{ID: 2, Email: "employee@test.com"}

	t.Run("Admin Reset Success", func(t *testing.T) {
		mockPasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return adminToken, nil
		}
		mockPasswordResetRepo.MarkPasswordResetTokenAsUsedFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}
		mockAdminCompanyRepo.GetAdminCompanyByIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mockAdminCompanyRepo.UpdateAdminCompanyFunc = func(ac *models.AdminCompaniesTable) error {
			return nil
		}

		err := service.ResetPassword("admin-token", "NewPassword123")

		assert.NoError(t, err)
	})

	t.Run("Employee Reset Success", func(t *testing.T) {
		mockPasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return employeeToken, nil
		}
		mockPasswordResetRepo.MarkPasswordResetTokenAsUsedFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockEmployeeRepo.UpdateEmployeeFunc = func(e *models.EmployeesTable) error {
			return nil
		}
		mockEmployeeRepo.SetEmployeePasswordSetFunc = func(id uint, isSet bool) error {
			return nil
		}

		err := service.ResetPassword("employee-token", "NewPassword123")

		assert.NoError(t, err)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockPasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return nil, errors.New("not found")
		}

		err := service.ResetPassword("invalid-token", "NewPassword123")

		assert.Error(t, err)
		assert.Equal(t, "failed to reset password: not found", err.Error())
	})
}
