package services_test

import (
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetMyLeaveRequests(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLeaveRequestService(nil, mocks.LeaveRequestRepo, nil)

	leaveRequests := []models.LeaveRequest{{Model: gorm.Model{ID: 1}}}

	t.Run("Success", func(t *testing.T) {
		mocks.LeaveRequestRepo.GetLeaveRequestsByEmployeeIDFunc = func(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
			return leaveRequests, nil
		}

		result, err := service.GetMyLeaveRequests(1, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, leaveRequests, result)
	})
}

func TestGetAllCompanyLeaveRequests(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLeaveRequestService(nil, mocks.LeaveRequestRepo, nil)

	leaveRequests := []models.LeaveRequest{{Model: gorm.Model{ID: 1}}}

	t.Run("Success", func(t *testing.T) {
		mocks.LeaveRequestRepo.GetCompanyLeaveRequestsPaginatedFunc = func(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error) {
			return leaveRequests, 1, nil
		}

		result, count, err := service.GetAllCompanyLeaveRequests(1, "", "", nil, nil, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, leaveRequests, result)
		assert.Equal(t, int64(1), count)
	})
}

func TestReviewLeaveRequest(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLeaveRequestService(mocks.EmployeeRepo, mocks.LeaveRequestRepo, mocks.AdminCompanyRepo)

	leaveRequest := &models.LeaveRequest{Model: gorm.Model{ID: 1}, EmployeeID: 1, Status: "pending"}
	employee := &models.EmployeesTable{ID: 1, CompanyID: 1}
	admin := &models.AdminCompaniesTable{ID: 1, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mocks.LeaveRequestRepo.GetLeaveRequestByIDFunc = func(id uint) (*models.LeaveRequest, error) {
			return leaveRequest, nil
		}
		mocks.EmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.AdminCompanyRepo.GetAdminCompanyByIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mocks.LeaveRequestRepo.UpdateLeaveRequestFunc = func(lr *models.LeaveRequest) error {
			return nil
		}

		result, err := service.ReviewLeaveRequest(1, 1, "approved")

		assert.NoError(t, err)
		assert.Equal(t, "approved", result.Status)
	})
}

func TestExportCompanyLeaveRequests(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLeaveRequestService(nil, mocks.LeaveRequestRepo, nil)

	leaveRequests := []models.LeaveRequest{{Model: gorm.Model{ID: 1}}}

	t.Run("Success", func(t *testing.T) {
		mocks.LeaveRequestRepo.GetCompanyLeaveRequestsFilteredFunc = func(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
			return leaveRequests, nil
		}

		result, err := service.ExportCompanyLeaveRequests(1, "", "", nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, leaveRequests, result)
	})
}

func TestCancelLeaveRequest(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLeaveRequestService(nil, mocks.LeaveRequestRepo, nil)

	leaveRequest := &models.LeaveRequest{Model: gorm.Model{ID: 1}, EmployeeID: 1, Status: "pending"}

	t.Run("Success", func(t *testing.T) {
		mocks.LeaveRequestRepo.GetLeaveRequestByIDFunc = func(id uint) (*models.LeaveRequest, error) {
			return leaveRequest, nil
		}
		mocks.LeaveRequestRepo.UpdateLeaveRequestFunc = func(lr *models.LeaveRequest) error {
			return nil
		}

		result, err := service.CancelLeaveRequest(1, 1)

		assert.NoError(t, err)
		assert.Equal(t, "cancelled", result.Status)
	})
}

func TestAdminCancelApprovedLeave(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewLeaveRequestService(mocks.EmployeeRepo, mocks.LeaveRequestRepo, mocks.AdminCompanyRepo)

	leaveRequest := &models.LeaveRequest{Model: gorm.Model{ID: 1}, EmployeeID: 1, Status: "approved"}
	employee := &models.EmployeesTable{ID: 1, CompanyID: 1}
	admin := &models.AdminCompaniesTable{ID: 1, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mocks.LeaveRequestRepo.GetLeaveRequestByIDFunc = func(id uint) (*models.LeaveRequest, error) {
			return leaveRequest, nil
		}
		mocks.EmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.AdminCompanyRepo.GetAdminCompanyByIDFunc = func(id int) (*models.AdminCompaniesTable, error) {
			return admin, nil
		}
		mocks.LeaveRequestRepo.UpdateLeaveRequestFunc = func(lr *models.LeaveRequest) error {
			return nil
		}

		result, err := service.AdminCancelApprovedLeave(1, 1)

		assert.NoError(t, err)
		assert.Equal(t, "cancelled", result.Status)
	})
}
