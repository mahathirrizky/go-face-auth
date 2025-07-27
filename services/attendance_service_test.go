package services

import (
	"errors"
	"go-face-auth/models"
	"testing"
	

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)



// MockFaceImageRepository is a mock implementation of FaceImageRepository for testing.
type MockFaceImageRepository struct {
	mock.Mock
}

func (m *MockFaceImageRepository) CreateFaceImage(faceImage *models.FaceImagesTable) error {
	panic("unimplemented")
}

func (m *MockFaceImageRepository) GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error) {
	args := m.Called(employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.FaceImagesTable), args.Error(1)
}

func (m *MockFaceImageRepository) GetFaceImageByID(id int) (*models.FaceImagesTable, error) {
	panic("unimplemented")
}

func (m *MockFaceImageRepository) DeleteFaceImage(id int) error {
	panic("unimplemented")
}

// MockAttendanceLocationRepository is a mock implementation of AttendanceLocationRepository for testing.
type MockAttendanceLocationRepository struct {
	mock.Mock
}

func (m *MockAttendanceLocationRepository) CreateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	panic("unimplemented")
}

func (m *MockAttendanceLocationRepository) GetAttendanceLocationsByCompanyID(companyID uint) ([]models.AttendanceLocation, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AttendanceLocation), args.Error(1)
}

func (m *MockAttendanceLocationRepository) GetAttendanceLocationByID(locationID uint) (*models.AttendanceLocation, error) {
	panic("unimplemented")
}

func (m *MockAttendanceLocationRepository) UpdateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	panic("unimplemented")
}

func (m *MockAttendanceLocationRepository) DeleteAttendanceLocation(locationID uint) error {
	panic("unimplemented")
}

// MockShiftRepository is a mock implementation of ShiftRepository for testing.
type MockShiftRepository struct {
	mock.Mock
}

func (m *MockShiftRepository) CreateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	panic("unimplemented")
}

func (m *MockShiftRepository) GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ShiftsTable), args.Error(1)
}

func (m *MockShiftRepository) GetShiftByID(id int) (*models.ShiftsTable, error) {
	panic("unimplemented")
}

func (m *MockShiftRepository) UpdateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	panic("unimplemented")
}

func (m *MockShiftRepository) DeleteShift(id int) error {
	panic("unimplemented")
}

func (m *MockShiftRepository) SetDefaultShift(companyID, shiftID int) error {
	panic("unimplemented")
}

func (m *MockShiftRepository) GetDefaultShiftByCompanyID(companyID int) (*models.ShiftsTable, error) {
	panic("unimplemented")
}

func TestHandleAttendance(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockCompanyRepo := new(MockCompanyRepository)
	mockAttendanceRepo := new(MockAttendanceRepository)
	mockFaceImageRepo := new(MockFaceImageRepository)
	mockLocationRepo := new(MockAttendanceLocationRepository)
	mockLeaveRequestRepo := new(MockLeaveRequestRepository)
	mockShiftRepo := new(MockShiftRepository)

	service := NewAttendanceService(mockEmployeeRepo, mockCompanyRepo, mockAttendanceRepo, mockFaceImageRepo, mockLocationRepo, mockLeaveRequestRepo, mockShiftRepo)

	req := AttendanceRequest{
		EmployeeID: 1,
		Latitude:   -6.200000,
		Longitude:  106.816666,
		ImageData:  "base64encodedimage",
	}

	shiftID := 1
	employee := &models.EmployeesTable{ID: 1, CompanyID: 1, ShiftID: &shiftID, Shift: models.ShiftsTable{ID: 1, StartTime: "09:00", EndTime: "17:00", GracePeriodMinutes: 15}}
	company := &models.CompaniesTable{ID: 1, Timezone: "Asia/Jakarta"}

	t.Run("Employee Not Found", func(t *testing.T) {
		mockEmployeeRepo.On("GetEmployeeByID", req.EmployeeID).Return(nil, errors.New("not found")).Once()

		_, _, _, err := service.HandleAttendance(req)

		assert.Error(t, err)
		assert.Equal(t, "employee not found", err.Error())
		mockEmployeeRepo.AssertExpectations(t)
	})

	t.Run("Company Not Found", func(t *testing.T) {
		mockEmployeeRepo.On("GetEmployeeByID", req.EmployeeID).Return(employee, nil).Once()
		mockCompanyRepo.On("GetCompanyByID", employee.CompanyID).Return(nil, errors.New("not found")).Once()

		_, _, _, err := service.HandleAttendance(req)

		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve company information", err.Error())
		mockEmployeeRepo.AssertExpectations(t)
		mockCompanyRepo.AssertExpectations(t)
	})

	t.Run("On Approved Leave", func(t *testing.T) {
		mockEmployeeRepo.On("GetEmployeeByID", req.EmployeeID).Return(employee, nil).Once()
		mockCompanyRepo.On("GetCompanyByID", employee.CompanyID).Return(company, nil).Once()
		mockLeaveRequestRepo.On("IsEmployeeOnApprovedLeave", employee.ID, mock.AnythingOfType("time.Time")).Return(&models.LeaveRequest{Type: "cuti"}, nil).Once()

		_, _, _, err := service.HandleAttendance(req)

		assert.Error(t, err)
		assert.Equal(t, "anda sedang dalam pengajuan cuti yang disetujui untuk hari ini", err.Error())
		mockEmployeeRepo.AssertExpectations(t)
		mockCompanyRepo.AssertExpectations(t)
		mockLeaveRequestRepo.AssertExpectations(t)
	})
}