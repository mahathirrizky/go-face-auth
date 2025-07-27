
package services

import (
	"errors"
	"go-face-auth/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	
)

// MockAdminCompanyRepository is a mock implementation of AdminCompanyRepository for testing.
type MockAdminCompanyRepository struct {
	mock.Mock
}

func (m *MockAdminCompanyRepository) CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	args := m.Called(adminCompany)
	return args.Error(0)
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AdminCompaniesTable), args.Error(1)
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error) {
	args := m.Called(employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AdminCompaniesTable), args.Error(1)
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByEmail(email string) (*models.AdminCompaniesTable, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AdminCompaniesTable), args.Error(1)
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByID(adminID int) (*models.AdminCompaniesTable, error) {
	args := m.Called(adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AdminCompaniesTable), args.Error(1)
}

func (m *MockAdminCompanyRepository) ChangeAdminPassword(adminID int, newPassword string) error {
	args := m.Called(adminID, newPassword)
	return args.Error(0)
}

func (m *MockAdminCompanyRepository) UpdateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	panic("unimplemented")
}

// MockCompanyRepository is a mock implementation of CompanyRepository for testing.
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) CreateCompany(company *models.CompaniesTable) error {
	panic("unimplemented")
}

func (m *MockCompanyRepository) GetCompanyByID(id int) (*models.CompaniesTable, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompaniesTable), args.Error(1)
}

func (m *MockCompanyRepository) GetCompanyWithSubscriptionDetails(id int) (*models.CompaniesTable, error) {
	panic("unimplemented")
}

func (m *MockCompanyRepository) UpdateCompany(company *models.CompaniesTable) error {
	args := m.Called(company)
	return args.Error(0)
}

func (m *MockCompanyRepository) GetAllActiveCompanies() ([]models.CompaniesTable, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CompaniesTable), args.Error(1)
}

// MockEmployeeRepository is a mock implementation of EmployeeRepository for testing.
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) GetTotalEmployeesByCompanyID(companyID int) (int64, error) {
	args := m.Called(companyID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEmployeeRepository) GetAbsentEmployeesCountToday(companyID int, presentToday int64) (int64, error) {
	args := m.Called(companyID, presentToday)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEmployeeRepository) CreateEmployee(employee *models.EmployeesTable) error {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetEmployeeByID(id int) (*models.EmployeesTable, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmployeesTable), args.Error(1)
}

func (m *MockEmployeeRepository) GetEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetEmployeeByEmail(email string) (*models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) UpdateEmployee(employee *models.EmployeesTable) error {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) DeleteEmployee(id int) error {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) SearchEmployees(companyID int, name string) ([]models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetEmployeesWithFaceImages(companyID int) ([]models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetOnLeaveEmployeesCountToday(companyID int) (int64, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) SetEmployeePasswordSet(employeeID uint, isSet bool) error {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) UpdateEmployeePassword(employee *models.EmployeesTable, newPassword string) error {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetPendingEmployees(companyID int) ([]models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetEmployeeByEmailOrIDNumber(email, employeeIDNumber string) (*models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetPendingEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) UpdateEmployeeFields(employee *models.EmployeesTable, updates map[string]interface{}) error {
	panic("unimplemented")
}

func (m *MockEmployeeRepository) GetActiveEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error) {
	panic("unimplemented")
}

// MockAttendanceRepository is a mock implementation of AttendanceRepository for testing.
type MockAttendanceRepository struct {
	mock.Mock
}

func (m *MockAttendanceRepository) GetPresentEmployeesCountToday(companyID int) (int64, error) {
	args := m.Called(companyID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAttendanceRepository) GetRecentAttendancesByCompanyID(companyID, limit int) ([]models.AttendancesTable, error) {
	args := m.Called(companyID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AttendancesTable), args.Error(1)
}

func (m *MockAttendanceRepository) GetRecentOvertimeAttendancesByCompanyID(companyID, limit int) ([]models.AttendancesTable, error) {
	args := m.Called(companyID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AttendancesTable), args.Error(1)
}

func (m *MockAttendanceRepository) CreateAttendance(attendance *models.AttendancesTable) error {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) UpdateAttendance(attendance *models.AttendancesTable) error {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetLatestAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetLatestAttendanceForDate(employeeID int, date time.Time) (*models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetLatestOvertimeAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetAbsentEmployeesCountToday(companyID int) (int64, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetAttendancesByCompanyID(companyID int) ([]models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetCompanyAttendancesFiltered(companyID int, startDate, endDate *time.Time, attendanceType string) ([]models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) HasAttendanceForDate(employeeID int, date time.Time) (bool, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) HasAttendanceForDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetCompanyOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetTodayAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetRecentAttendancesByEmployeeID(employeeID int, limit int) ([]models.AttendancesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetOvertimeAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetUnaccountedEmployeesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetUnaccountedEmployeesFiltered(companyID int, startDate, endDate *time.Time, search string) ([]models.EmployeesTable, error) {
	panic("unimplemented")
}

func (m *MockAttendanceRepository) GetOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time, search string) ([]models.AttendancesTable, error) {
	panic("unimplemented")
}

// MockLeaveRequestRepository is a mock implementation of LeaveRequestRepository for testing.
type MockLeaveRequestRepository struct {
	mock.Mock
}

func (m *MockLeaveRequestRepository) GetOnLeaveEmployeesCountToday(companyID int) (int64, error) {
	args := m.Called(companyID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLeaveRequestRepository) GetRecentLeaveRequestsByCompanyID(companyID, limit int) ([]models.LeaveRequest, error) {
	args := m.Called(companyID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LeaveRequest), args.Error(1)
}

func (m *MockLeaveRequestRepository) CreateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) GetLeaveRequestByID(id uint) (*models.LeaveRequest, error) {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) GetAllLeaveRequests() ([]models.LeaveRequest, error) {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) GetLeaveRequestsByEmployeeID(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) GetCompanyLeaveRequestsFiltered(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) UpdateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) IsEmployeeOnApprovedLeave(employeeID int, date time.Time) (*models.LeaveRequest, error) {
	args := m.Called(employeeID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LeaveRequest), args.Error(1)
}

func (m *MockLeaveRequestRepository) IsEmployeeOnApprovedLeaveDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) GetPendingLeaveRequestsByEmployeeID(employeeID int) ([]models.LeaveRequest, error) {
	panic("unimplemented")
}

func (m *MockLeaveRequestRepository) GetCompanyLeaveRequestsPaginated(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error) {
	panic("unimplemented")
}

func TestCreateAdminCompany(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil, nil)

	admin := &models.AdminCompaniesTable{Email: "test@example.com"}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.On("CreateAdminCompany", admin).Return(nil).Once()

		err := service.CreateAdminCompany(admin)

		assert.NoError(t, err)
		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("failed to create admin")
		mockAdminRepo.On("CreateAdminCompany", admin).Return(expectedError).Once()

		err := service.CreateAdminCompany(admin)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockAdminRepo.AssertExpectations(t)
	})
}

func TestGetAdminCompanyByCompanyID(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil, nil)

	companyID := 1
	admin := &models.AdminCompaniesTable{ID: 1, CompanyID: companyID}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.On("GetAdminCompanyByCompanyID", companyID).Return(admin, nil).Once()

		result, err := service.GetAdminCompanyByCompanyID(companyID)

		assert.NoError(t, err)
		assert.Equal(t, admin, result)
		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("not found")
		mockAdminRepo.On("GetAdminCompanyByCompanyID", companyID).Return(nil, expectedError).Once()

		result, err := service.GetAdminCompanyByCompanyID(companyID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockAdminRepo.AssertExpectations(t)
	})
}

func TestGetAdminCompanyByEmployeeID(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil, nil)

	employeeID := 1
	admin := &models.AdminCompaniesTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.On("GetAdminCompanyByEmployeeID", employeeID).Return(admin, nil).Once()

		result, err := service.GetAdminCompanyByEmployeeID(employeeID)

		assert.NoError(t, err)
		assert.Equal(t, admin, result)
		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("not found")
		mockAdminRepo.On("GetAdminCompanyByEmployeeID", employeeID).Return(nil, expectedError).Once()

		result, err := service.GetAdminCompanyByEmployeeID(employeeID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockAdminRepo.AssertExpectations(t)
	})
}

func TestChangeAdminPassword(t *testing.T) {
	mockAdminRepo := new(MockAdminCompanyRepository)
	service := NewAdminCompanyService(mockAdminRepo, nil, nil, nil, nil, nil)

	adminID := 1
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	admin := &models.AdminCompaniesTable{ID: 1, Password: string(hashedOldPassword)}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.On("GetAdminCompanyByID", adminID).Return(admin, nil).Once()
		mockAdminRepo.On("ChangeAdminPassword", adminID, mock.AnythingOfType("string")).Return(nil).Once()

		err := service.ChangeAdminPassword(adminID, oldPassword, newPassword)

		assert.NoError(t, err)
		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Admin Not Found", func(t *testing.T) {
		mockAdminRepo.On("GetAdminCompanyByID", adminID).Return(nil, errors.New("not found")).Once()

		err := service.ChangeAdminPassword(adminID, oldPassword, newPassword)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve admin details")
		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Incorrect Old Password", func(t *testing.T) {
		mockAdminRepo.On("GetAdminCompanyByID", adminID).Return(admin, nil).Once()

		err := service.ChangeAdminPassword(adminID, "wrongPassword", newPassword)

		assert.Error(t, err)
		assert.Equal(t, "incorrect old password", err.Error())
		mockAdminRepo.AssertExpectations(t)
	})
}

func TestGetDashboardSummaryData(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockAttendanceRepo := new(MockAttendanceRepository)
	mockLeaveRepo := new(MockLeaveRequestRepository)
	service := NewAdminCompanyService(nil, nil, mockEmployeeRepo, mockAttendanceRepo, mockLeaveRepo, nil)

	companyID := 1

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.On("GetTotalEmployeesByCompanyID", companyID).Return(int64(10), nil).Once()
		mockAttendanceRepo.On("GetPresentEmployeesCountToday", companyID).Return(int64(8), nil).Once()
		mockEmployeeRepo.On("GetAbsentEmployeesCountToday", companyID, int64(8)).Return(int64(2), nil).Once()
		mockLeaveRepo.On("GetOnLeaveEmployeesCountToday", companyID).Return(int64(1), nil).Once()
		mockAttendanceRepo.On("GetRecentAttendancesByCompanyID", companyID, 10).Return([]models.AttendancesTable{}, nil).Once()
		mockLeaveRepo.On("GetRecentLeaveRequestsByCompanyID", companyID, 10).Return([]models.LeaveRequest{}, nil).Once()
		mockAttendanceRepo.On("GetRecentOvertimeAttendancesByCompanyID", companyID, 10).Return([]models.AttendancesTable{}, nil).Once()

		summary, err := service.GetDashboardSummaryData(companyID)

		assert.NoError(t, err)
		assert.Equal(t, int64(10), summary["total_employees"])
		assert.Equal(t, int64(8), summary["present_today"])
		assert.Equal(t, int64(2), summary["absent_today"])
		assert.Equal(t, int64(1), summary["on_leave_today"])
		mockEmployeeRepo.AssertExpectations(t)
		mockAttendanceRepo.AssertExpectations(t)
		mockLeaveRepo.AssertExpectations(t)
	})

	t.Run("Error Fetching Total Employees", func(t *testing.T) {
		expectedError := errors.New("db error")
		mockEmployeeRepo.On("GetTotalEmployeesByCompanyID", companyID).Return(int64(0), expectedError).Once()

		summary, err := service.GetDashboardSummaryData(companyID)

		assert.Error(t, err)
		assert.Nil(t, summary)
		assert.Equal(t, expectedError, err)
		mockEmployeeRepo.AssertExpectations(t)
	})
}
