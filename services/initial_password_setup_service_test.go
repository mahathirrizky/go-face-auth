package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetupInitialPassword(t *testing.T) {
	mockPasswordResetRepo := new(MockPasswordResetRepository)
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewInitialPasswordSetupService(mockPasswordResetRepo, mockEmployeeRepo)

	token := &models.PasswordResetTokenTable{Token: "test-token", UserID: 1, ExpiresAt: time.Now().Add(time.Hour), TokenType: "employee_initial_password"}
	employee := &models.EmployeesTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mockPasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return token, nil
		}
		mockPasswordResetRepo.MarkPasswordResetTokenAsUsedFunc = func(t *models.PasswordResetTokenTable) error {
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

		err := service.SetupInitialPassword("test-token", "newPassword1A")

		assert.NoError(t, err)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		err := service.SetupInitialPassword("test-token", "short")

		assert.Error(t, err)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockPasswordResetRepo.GetPasswordResetTokenFunc = func(tokenString string) (*models.PasswordResetTokenTable, error) {
			return nil, errors.New("not found")
		}

		err := service.SetupInitialPassword("invalid-token", "newPassword1A")

		assert.Error(t, err)
	})
}
