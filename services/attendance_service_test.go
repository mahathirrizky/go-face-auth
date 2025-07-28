package services_test

import (
	"errors"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockPythonServerClient is a mock of PythonServerClientInterface
type MockPythonServerClient struct {
	SendToPythonServerFunc func(payload services.PythonRecognitionRequest) (map[string]interface{}, error)
}

// SendToPythonServer mocks the SendToPythonServer method
func (m *MockPythonServerClient) SendToPythonServer(payload services.PythonRecognitionRequest) (map[string]interface{}, error) {
	if m.SendToPythonServerFunc != nil {
		return m.SendToPythonServerFunc(payload)
	}
	return nil, errors.New("SendToPythonServerFunc not implemented")
}

func TestHandleAttendance(t *testing.T) {
	mocks := services.NewMockRepositories()
	mockPythonClient := new(MockPythonServerClient)

	service := services.NewAttendanceService(mocks.EmployeeRepo, mocks.CompanyRepo, mocks.AttendanceRepo, mocks.FaceImageRepo, mocks.AttendanceLocationRepo, mocks.LeaveRequestRepo, mocks.ShiftRepo, mocks.DivisionRepo, mockPythonClient)

	req := services.AttendanceRequest{
		EmployeeID: 1,
		Latitude:   -6.200000,
		Longitude:  106.816666,
		ImageData:  "base64encodedimage",
	}

	shiftID := 1
	employee := &models.EmployeesTable{ID: 1, CompanyID: 1, ShiftID: &shiftID, Shift: models.ShiftsTable{ID: 1, StartTime: "09:00", EndTime: "17:00", GracePeriodMinutes: 15}}
	company := &models.CompaniesTable{ID: 1, Timezone: "Asia/Jakarta"}

	t.Run("Employee Not Found", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return nil, errors.New("not found")
		}

		_, _, _, err := service.HandleAttendance(req)

		assert.Error(t, err)
		assert.Equal(t, "employee not found", err.Error())
	})

	t.Run("Company Not Found", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.CompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return nil, errors.New("not found")
		}

		_, _, _, err := service.HandleAttendance(req)

		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve company information", err.Error())
	})

	t.Run("On Approved Leave", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.CompanyRepo.GetCompanyByIDFunc = func(id int) (*models.CompaniesTable, error) {
			return company, nil
		}
		mocks.LeaveRequestRepo.IsEmployeeOnApprovedLeaveFunc = func(employeeID int, date time.Time) (*models.LeaveRequest, error) {
			return &models.LeaveRequest{Type: "cuti"}, nil
		}

		_, _, _, err := service.HandleAttendance(req)

		assert.Error(t, err)
		assert.Equal(t, "anda sedang dalam pengajuan cuti yang disetujui untuk hari ini", err.Error())
	})
}
