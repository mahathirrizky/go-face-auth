package repository

import "go-face-auth/models"

// EmployeeRepository defines the contract for employee-related database operations.
type EmployeeRepository interface {
	CreateEmployee(employee *models.EmployeesTable) error
	GetEmployeeByID(id int) (*models.EmployeesTable, error)
	GetEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error)
	GetEmployeeByEmail(email string) (*models.EmployeesTable, error)
	UpdateEmployee(employee *models.EmployeesTable) error
	DeleteEmployee(id int) error
	SearchEmployees(companyID int, name string) ([]models.EmployeesTable, error)
	GetTotalEmployeesByCompanyID(companyID int) (int64, error)
	GetEmployeesWithFaceImages(companyID int) ([]models.EmployeesTable, error)
	GetOnLeaveEmployeesCountToday(companyID int) (int64, error)
	SetEmployeePasswordSet(employeeID uint, isSet bool) error
	UpdateEmployeePassword(employee *models.EmployeesTable, newPassword string) error
	GetPendingEmployees(companyID int) ([]models.EmployeesTable, error)
	GetEmployeeByEmailOrIDNumber(email, employeeIDNumber string) (*models.EmployeesTable, error)
	GetEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error)
	GetPendingEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error)
	UpdateEmployeeFields(employee *models.EmployeesTable, updates map[string]interface{}) error
	GetActiveEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error)
	// This function seems to be missing from the original repo but was in the service, let's define it.
	GetAbsentEmployeesCountToday(companyID int, presentToday int64) (int64, error)
}
