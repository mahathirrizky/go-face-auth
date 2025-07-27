package services_test

import (
	"errors"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmployeeByID(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employee := &models.EmployeesTable{ID: 1, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}

		result, err := service.GetEmployeeByID(1, 1)

		assert.NoError(t, err)
		assert.Equal(t, employee, result)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.GetEmployeeByID(1, 1)

		assert.Error(t, err)
	})

	t.Run("Wrong Company", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}

		result, err := service.GetEmployeeByID(1, 2)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestGetEmployeesByCompanyIDPaginated(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employees := []models.EmployeesTable{{ID: 1}, {ID: 2}}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeesByCompanyIDPaginatedFunc = func(companyID int, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error) {
			return employees, 2, nil
		}

		result, count, err := service.GetEmployeesByCompanyIDPaginated(1, "", 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, employees, result)
		assert.Equal(t, int64(2), count)
	})
}

func TestSearchEmployees(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employees := []models.EmployeesTable{{Name: "John Doe"}}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.SearchEmployeesFunc = func(companyID int, name string) ([]models.EmployeesTable, error) {
			return employees, nil
		}

		result, err := service.SearchEmployees(1, "John")

		assert.NoError(t, err)
		assert.Equal(t, employees, result)
	})
}

func TestUpdateEmployee(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employee := &models.EmployeesTable{ID: 1, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockEmployeeRepo.UpdateEmployeeFieldsFunc = func(e *models.EmployeesTable, updates map[string]interface{}) error {
			return nil
		}

		err := service.UpdateEmployee(1, 1, map[string]interface{}{"name": "New Name"})

		assert.NoError(t, err)
	})
}

func TestDeleteEmployee(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employee := &models.EmployeesTable{ID: 1, CompanyID: 1}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockEmployeeRepo.DeleteEmployeeFunc = func(id int) error {
			return nil
		}

		err := service.DeleteEmployee(1, 1)

		assert.NoError(t, err)
	})
}

func TestGetPendingEmployeesByCompanyIDPaginated(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employees := []models.EmployeesTable{{IsPasswordSet: false}}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetPendingEmployeesByCompanyIDPaginatedFunc = func(companyID int, search string, page int, pageSize int) ([]models.EmployeesTable, int64, error) {
			return employees, 1, nil
		}

		result, count, err := service.GetPendingEmployeesByCompanyIDPaginated(1, "", 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, employees, result)
		assert.Equal(t, int64(1), count)
	})
}

func TestUpdateEmployeeProfile(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	employee := &models.EmployeesTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockEmployeeRepo.UpdateEmployeeFunc = func(e *models.EmployeesTable) error {
			return nil
		}

		req := services.UpdateEmployeeProfileRequest{Name: "New Name", Email: "new@email.com", Position: "New Pos"}
		err := service.UpdateEmployeeProfile(1, req)

		assert.NoError(t, err)
	})
}

func TestChangeEmployeePassword(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, nil, nil, nil)

	password, _ := helper.HashPassword("oldPassword")
	employee := &models.EmployeesTable{ID: 1, Password: password}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockEmployeeRepo.UpdateEmployeePasswordFunc = func(e *models.EmployeesTable, newPassword string) error {
			return nil
		}

		err := service.ChangeEmployeePassword(1, "oldPassword", "newPassword1A", "newPassword1A")

		assert.NoError(t, err)
	})
}

func TestGetEmployeeDashboardSummary(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockAttendanceRepo := new(MockAttendanceRepository)
	mockLeaveRequestRepo := new(MockLeaveRequestRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, nil, nil, nil, mockAttendanceRepo, mockLeaveRequestRepo, nil)

	employee := &models.EmployeesTable{ID: 1}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockAttendanceRepo.GetTodayAttendanceByEmployeeIDFunc = func(id int) (*models.AttendancesTable, error) {
			return &models.AttendancesTable{Status: "Checked In"}, nil
		}
		mockLeaveRequestRepo.GetPendingLeaveRequestsByEmployeeIDFunc = func(id int) ([]models.LeaveRequest, error) {
			return []models.LeaveRequest{}, nil
		}
		mockAttendanceRepo.GetRecentAttendancesByEmployeeIDFunc = func(id int, limit int) ([]models.AttendancesTable, error) {
			return []models.AttendancesTable{}, nil
		}

		summary, err := service.GetEmployeeDashboardSummary(1)

		assert.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, "Checked In", summary.TodayAttendanceStatus)
	})
}

func TestGetEmployeeProfile(t *testing.T) {
	mockEmployeeRepo := new(MockEmployeeRepository)
	mockShiftRepo := new(MockShiftRepository)
	mockAttendanceLocationRepo := new(MockAttendanceLocationRepository)
	mockFaceImageRepo := new(MockFaceImageRepository)
	service := services.NewEmployeeService(mockEmployeeRepo, nil, mockShiftRepo, nil, mockFaceImageRepo, nil, nil, mockAttendanceLocationRepo)

	shiftID := 1
	employee := &models.EmployeesTable{ID: 1, ShiftID: &shiftID}

	t.Run("Success", func(t *testing.T) {
		mockEmployeeRepo.GetEmployeeByIDFunc = func(id int) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mockShiftRepo.GetShiftByIDFunc = func(id int) (*models.ShiftsTable, error) {
			return &models.ShiftsTable{ID: 1}, nil
		}
		mockAttendanceLocationRepo.GetAttendanceLocationsByCompanyIDFunc = func(id uint) ([]models.AttendanceLocation, error) {
			return []models.AttendanceLocation{}, nil
		}
		mockFaceImageRepo.GetFaceImagesByEmployeeIDFunc = func(id int) ([]models.FaceImagesTable, error) {
			return []models.FaceImagesTable{}, nil
		}

		profile, err := service.GetEmployeeProfile(1)

		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, employee.ID, profile.ID)
	})
}
