
package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAdminCompany(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := services.NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil)

	admin := &models.AdminCompaniesTable{Email: "test@example.com"}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.CreateAdminCompanyFunc = func(adminCompany *models.AdminCompaniesTable) error {
			return nil
		}

		err := service.CreateAdminCompany(admin)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("failed to create admin")
		mockAdminRepo.CreateAdminCompanyFunc = func(adminCompany *models.AdminCompaniesTable) error {
			return expectedError
		}

		err := service.CreateAdminCompany(admin)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestGetAdminCompanyByCompanyID(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := services.NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil)

	companyID := 1
	admin := &models.AdminCompaniesTable{ID: 1, CompanyID: companyID}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.GetAdminCompanyByCompanyIDFunc = func(cID int) (*models.AdminCompaniesTable, error) {
			assert.Equal(t, companyID, cID)
			return admin, nil
		}

		result, err := service.GetAdminCompanyByCompanyID(companyID)

		assert.NoError(t, err)
		assert.Equal(t, admin, result)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("not found")
		mockAdminRepo.GetAdminCompanyByCompanyIDFunc = func(cID int) (*models.AdminCompaniesTable, error) {
			return nil, expectedError
		}

		result, err := service.GetAdminCompanyByCompanyID(companyID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestGetAdminCompanyByEmployeeID(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := services.NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil)

	employeeID := 1
	admin := &models.AdminCompaniesTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.GetAdminCompanyByEmployeeIDFunc = func(eID int) (*models.AdminCompaniesTable, error) {
			assert.Equal(t, employeeID, eID)
			return admin, nil
		}

		result, err := service.GetAdminCompanyByEmployeeID(employeeID)

		assert.NoError(t, err)
		assert.Equal(t, admin, result)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("not found")
		mockAdminRepo.GetAdminCompanyByEmployeeIDFunc = func(eID int) (*models.AdminCompaniesTable, error) {
			return nil, expectedError
		}

		result, err := service.GetAdminCompanyByEmployeeID(employeeID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestChangeAdminPassword(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := services.NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil)

	adminID := 1
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	admin := &models.AdminCompaniesTable{ID: 1, Password: string(hashedOldPassword)}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.GetAdminCompanyByIDFunc = func(aID int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
				mockAdminRepo.ChangeAdminPasswordFunc = func(aID uint, newPass string) error {
			return nil
		}

		err := service.ChangeAdminPassword(adminID, oldPassword, newPassword)

		assert.NoError(t, err)
	})

	t.Run("Admin Not Found", func(t *testing.T) {
		mockAdminRepo.GetAdminCompanyByIDFunc = func(aID int) (*models.AdminCompaniesTable, error) {
			return nil, errors.New("not found")
		}

		err := service.ChangeAdminPassword(adminID, oldPassword, newPassword)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve admin details")
	})

	t.Run("Incorrect Old Password", func(t *testing.T) {
		mockAdminRepo.GetAdminCompanyByIDFunc = func(aID int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}

		err := service.ChangeAdminPassword(adminID, "wrongPassword", newPassword)

		assert.Error(t, err)
		assert.Equal(t, "incorrect old password", err.Error())
	})
}

