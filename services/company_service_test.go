package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCompany(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	service := services.NewCompanyService(mockCompanyRepo, nil, nil, nil)

	req := services.CreateCompanyRequest{Name: "Test Company", Address: "123 Test St"}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.CreateCompanyFunc = func(company *models.CompaniesTable) error {
			return nil
		}

		result, err := service.CreateCompany(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("Error", func(t *testing.T) {
		mockCompanyRepo.CreateCompanyFunc = func(company *models.CompaniesTable) error {
			return errors.New("db error")
		}

		_, err := service.CreateCompany(req)

		assert.Error(t, err)
	})
}

func TestGetCompanyByID(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	service := services.NewCompanyService(mockCompanyRepo, nil, nil, nil)

	company := &models.CompaniesTable{ID: 1, Name: "Test Company"}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}

		result, err := service.GetCompanyByID(1)

		assert.NoError(t, err)
		assert.Equal(t, company, result)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.GetCompanyByID(1)

		assert.Error(t, err)
	})
}

func TestGetCompanyDetails(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockAdminCompanyRepo := new(MockAdminCompanyRepository)
	service := services.NewCompanyService(mockCompanyRepo, mockAdminCompanyRepo, nil, nil)

	company := &models.CompaniesTable{ID: 1, Name: "Test Company"}
	adminCompany := &models.AdminCompaniesTable{Email: "admin@test.com"}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockAdminCompanyRepo.GetAdminCompanyByCompanyIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return adminCompany, nil
		}

		details, err := service.GetCompanyDetails(1)

		assert.NoError(t, err)
		assert.Equal(t, company.Name, details["name"])
		assert.Equal(t, adminCompany.Email, details["admin_email"])
	})
}

func TestUpdateCompanyDetails(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	service := services.NewCompanyService(mockCompanyRepo, nil, nil, nil)

	company := &models.CompaniesTable{ID: 1, Name: "Old Name", Address: "Old Address", Timezone: "UTC"}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mockCompanyRepo.UpdateCompanyFunc = func(c *models.CompaniesTable) error {
			return nil
		}

		updatedCompany, err := service.UpdateCompanyDetails(1, "New Name", "New Address", "Asia/Jakarta")

		assert.NoError(t, err)
		assert.Equal(t, "New Name", updatedCompany.Name)
		assert.Equal(t, "New Address", updatedCompany.Address)
		assert.Equal(t, "Asia/Jakarta", updatedCompany.Timezone)
	})
}

func TestGetCompanySubscriptionStatus(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	service := services.NewCompanyService(mockCompanyRepo, nil, nil, nil)

	company := &models.CompaniesTable{ID: 1, SubscriptionStatus: "trial"}

	t.Run("Success", func(t *testing.T) {
		mockCompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}

		status, err := service.GetCompanySubscriptionStatus(1)

		assert.NoError(t, err)
		assert.Equal(t, company.SubscriptionStatus, status["subscription_status"])
	})
}
