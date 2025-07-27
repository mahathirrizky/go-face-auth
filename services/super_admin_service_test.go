package services_test

import (
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSuperAdminDashboardSummary(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	mockSuperAdminRepo := new(MockSuperAdminRepository)
	service := services.NewSuperAdminService(mockCompanyRepo, mockInvoiceRepo, mockCustomPackageRequestRepo, mockSuperAdminRepo)

	t.Run("Success", func(t *testing.T) {
		mockSuperAdminRepo.GetTotalCompaniesCountFunc = func() (int64, error) {
			return 10, nil
		}
		mockSuperAdminRepo.GetCompaniesCountBySubscriptionStatusFunc = func(status string) (int64, error) {
			switch status {
			case "active":
				return 5, nil
			case "trial":
				return 3, nil
			default:
				return 0, nil
			}
		}
		mockSuperAdminRepo.GetExpiredAndTrialExpiredCompaniesCountFunc = func() (int64, error) {
			return 2, nil
		}
		mockSuperAdminRepo.GetRecentCompaniesFunc = func(limit int) ([]models.CompaniesTable, error) {
			return []models.CompaniesTable{{ID: 1, Name: "Comp1", CreatedAt: time.Now()}}, nil
		}

		summary, err := service.GetSuperAdminDashboardSummary()

		assert.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, int64(10), summary["total_companies"])
		assert.Equal(t, int64(5), summary["active_subscriptions"])
		assert.Equal(t, int64(2), summary["expired_subscriptions"])
		assert.Equal(t, int64(3), summary["trial_subscriptions"])
		assert.Len(t, summary["recent_activities"].([]map[string]interface{}), 1)
	})
}

func TestGetCompanies(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	mockSuperAdminRepo := new(MockSuperAdminRepository)
	service := services.NewSuperAdminService(mockCompanyRepo, mockInvoiceRepo, mockCustomPackageRequestRepo, mockSuperAdminRepo)

	companies := []models.CompaniesTable{{ID: 1, Name: "Comp1"}}

	t.Run("Success", func(t *testing.T) {
		mockSuperAdminRepo.GetAllCompaniesWithPreloadFunc = func() ([]models.CompaniesTable, error) {
			return companies, nil
		}

		result, err := service.GetCompanies()

		assert.NoError(t, err)
		assert.Equal(t, companies, result)
	})
}

func TestGetSubscriptions(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	mockSuperAdminRepo := new(MockSuperAdminRepository)
	service := services.NewSuperAdminService(mockCompanyRepo, mockInvoiceRepo, mockCustomPackageRequestRepo, mockSuperAdminRepo)

	subscriptions := []models.CompaniesTable{{ID: 1, Name: "SubComp1"}}

	t.Run("Success", func(t *testing.T) {
		mockSuperAdminRepo.GetAllCompaniesWithPreloadFunc = func() ([]models.CompaniesTable, error) {
			return subscriptions, nil
		}

		result, err := service.GetSubscriptions()

		assert.NoError(t, err)
		assert.Equal(t, subscriptions, result)
	})
}

func TestGetRevenueSummary(t *testing.T) {
	mockCompanyRepo := new(MockCompanyRepository)
	mockInvoiceRepo := new(MockInvoiceRepository)
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	mockSuperAdminRepo := new(MockSuperAdminRepository)
	service := services.NewSuperAdminService(mockCompanyRepo, mockInvoiceRepo, mockCustomPackageRequestRepo, mockSuperAdminRepo)

	revenue := []struct {
		Month        string
		Year         string
		TotalRevenue float64
	}{{"07", "2025", 1000.0}}

	t.Run("Success", func(t *testing.T) {
		mockSuperAdminRepo.GetPaidInvoicesMonthlyRevenueFunc = func(startDate, endDate *time.Time) ([]struct {
			Month        string
			Year         string
			TotalRevenue float64
		}, error) {
			return revenue, nil
		}

		result, err := service.GetRevenueSummary("", "")

		assert.NoError(t, err)
		assert.Equal(t, revenue[0].TotalRevenue, result[0].TotalRevenue)
	})
}

func TestGetCustomPackageRequests(t *testing.T) {
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	service := services.NewSuperAdminService(nil, nil, mockCustomPackageRequestRepo, nil)

	requests := []models.CustomPackageRequest{{ID: 1}}

	t.Run("Success", func(t *testing.T) {
		mockCustomPackageRequestRepo.GetCustomPackageRequestsPaginatedFunc = func(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error) {
			return requests, 1, nil
		}

		result, count, err := service.GetCustomPackageRequests(1, 10, "")

		assert.NoError(t, err)
		assert.Equal(t, requests, result)
		assert.Equal(t, int64(1), count)
	})
}

func TestUpdateCustomPackageRequestStatus(t *testing.T) {
	mockCustomPackageRequestRepo := new(MockCustomPackageRequestRepository)
	service := services.NewSuperAdminService(nil, nil, mockCustomPackageRequestRepo, nil)

	request := &models.CustomPackageRequest{ID: 1, Status: "pending"}

	t.Run("Success", func(t *testing.T) {
		mockCustomPackageRequestRepo.GetCustomPackageRequestByIDFunc = func(id uint) (*models.CustomPackageRequest, error) {
			return request, nil
		}
		mockCustomPackageRequestRepo.UpdateCustomPackageRequestFunc = func(req *models.CustomPackageRequest) error {
			return nil
		}

		err := service.UpdateCustomPackageRequestStatus(1, "approved")

		assert.NoError(t, err)
		assert.Equal(t, "approved", request.Status)
	})
}
