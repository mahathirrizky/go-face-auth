package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCustomPackageRequest(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockAdminCompanyRepo := new(MockAdminCompanyRepository)
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	service := services.NewCustomPackageRequestService(mockCompanyRepo, mockAdminCompanyRepo, mockCustomPackageRequestRepo)

	company := &models.CompaniesTable{ID: 1, Name: "Test Co"}
	admin := &models.AdminCompaniesTable{ID: 1, Email: "admin@test.com"}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockAdminCompanyRepo.GetAdminCompanyByIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mockCustomPackageRequestRepo.CreateCustomPackageRequestFunc = func(req *models.CustomPackageRequest) error {
			return nil
		}

		request, err := service.CreateCustomPackageRequest(1, 1, "12345", "test message")

		assert.NoError(t, err)
		assert.NotNil(t, request)
		assert.Equal(t, company.Name, request.CompanyName)
		assert.Equal(t, admin.Email, request.Email)
	})

	t.Run("Company Not Found", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.CreateCustomPackageRequest(1, 1, "12345", "test message")

		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve company information", err.Error())
	})

	t.Run("Admin Not Found", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockAdminCompanyRepo.GetAdminCompanyByIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.CreateCustomPackageRequest(1, 1, "12345", "test message")

		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve admin information", err.Error())
	})
}
