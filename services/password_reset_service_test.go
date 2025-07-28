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
	mocks := services.NewMockRepositories()
	service := services.NewPasswordResetService(mocks.AdminCompanyRepo, mocks.EmployeeRepo, mocks.PasswordResetRepo)

	admin := &models.AdminCompaniesTable{ID: 1, Email: "admin@test.com"}
	employee := &models.EmployeesTable{ID: 2, Email: "employee@test.com", Name: "Test Employee"}

	t.Run("Admin Success", func(t *testing.T) {
		mocks.AdminCompanyRepo.GetAdminCompanyByEmailFunc = func(email string) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mocks.PasswordResetRepo.CreatePasswordResetTokenFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}

		err := service.ForgotPassword("admin@test.com", "admin")

		assert.NoError(t, err)
	})

	t.Run("Employee Success", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByEmailFunc = func(email string) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.PasswordResetRepo.CreatePasswordResetTokenFunc = func(token *models.PasswordResetTokenTable) error {
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
	mocks := services.NewMockRepositories()
	service := services.NewPasswordResetService(mocks.AdminCompanyRepo, mocks.EmployeeRepo, mocks.PasswordResetRepo)

	adminToken := &models.PasswordResetTokenTable{Token: "admin-token", UserID: 1, ExpiresAt: time.Now().Add(time.Hour), TokenType: "admin_password_reset"}
	employeeToken := &models.PasswordResetTokenTable{Token: "employee-token", UserID: 2, ExpiresAt: time.Now().Add(time.Hour), TokenType: "employee_password_reset"}

	admin := &models.AdminCompaniesTable{ID: 1, Email: "admin@test.com"}
	employee := &models.EmployeesTable{ID: 2, Email: "employee@test.com"}

	t.Run("Admin Reset Success", func(t *testing.T) {
		mocks.PasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return adminToken, nil
		}
		mocks.PasswordResetRepo.MarkPasswordResetTokenAsUsedFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}
		mocks.AdminCompanyRepo.GetAdminCompanyByIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mocks.AdminCompanyRepo.UpdateAdminCompanyFunc = func(ac *models.AdminCompaniesTable) error {
			return nil
		}

		err := service.ResetPassword("admin-token", "NewPassword123")

		assert.NoError(t, err)
	})

	t.Run("Employee Reset Success", func(t *testing.T) {
		mocks.PasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return employeeToken, nil
		}
		mocks.PasswordResetRepo.MarkPasswordResetTokenAsUsedFunc = func(token *models.PasswordResetTokenTable) error {
			return nil
		}
		mocks.EmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.EmployeeRepo.UpdateEmployeeFunc = func(e *models.EmployeesTable) error {
			return nil
		}
		mocks.EmployeeRepo.SetEmployeePasswordSetFunc = func(id uint, isSet bool) error {
			return nil
		}

		err := service.ResetPassword("employee-token", "NewPassword123")

		assert.NoError(t, err)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mocks.PasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return nil, errors.New("not found")
		}

		err := service.ResetPassword("invalid-token", "NewPassword123")

		assert.Error(t, err)
		assert.Equal(t, "failed to reset password: not found", err.Error())
	})
}
