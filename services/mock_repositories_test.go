package services

import (
	"go-face-auth/database/repository"
	"go-face-auth/models"
	"time"

	"gorm.io/gorm"
	"reflect"
	"github.com/stretchr/testify/mock"
)

// GormDBInterface defines the methods of *gorm.DB used in services that need to be mocked.
type GormDBInterface interface {
	Model(value interface{}) GormDBInterface
	Where(query interface{}, args ...interface{}) GormDBInterface
	Count(count *int64) GormDBInterface
	Error() error
	Order(value interface{}) GormDBInterface
	Limit(limit int) GormDBInterface
	Find(out interface{}, where ...interface{}) GormDBInterface
	Preload(query string, args ...interface{}) GormDBInterface
	Select(query interface{}, args ...interface{}) GormDBInterface
	Group(name string) GormDBInterface
	Scan(dest interface{}) GormDBInterface
}

// MockGormDB is a mock for *gorm.DB for specific methods used in services tests
type MockGormDB struct {
	ErrorValue error
	CountValue int64
	FindValue  interface{}
	ScanValue  interface{}
}

func (m *MockGormDB) Model(value interface{}) GormDBInterface {
	return m // Return itself to allow chaining
}

func (m *MockGormDB) Where(query interface{}, args ...interface{}) GormDBInterface {
	return m // Return itself to allow chaining
}

func (m *MockGormDB) Count(count *int64) GormDBInterface {
	*count = m.CountValue
	return m // Return itself to allow chaining
}

func (m *MockGormDB) Order(value interface{}) GormDBInterface {
	return m
}

func (m *MockGormDB) Limit(limit int) GormDBInterface {
	return m
}

func (m *MockGormDB) Find(out interface{}, where ...interface{}) GormDBInterface {
	if m.FindValue != nil {
		srcVal := reflect.ValueOf(m.FindValue)
		dstVal := reflect.ValueOf(out)

		if srcVal.Kind() == reflect.Slice && dstVal.Kind() == reflect.Ptr && dstVal.Elem().Kind() == reflect.Slice {
			dstVal.Elem().Set(srcVal)
		} else if srcVal.Kind() == reflect.Ptr && dstVal.Kind() == reflect.Ptr && srcVal.Type() == dstVal.Type() {
			dstVal.Elem().Set(srcVal.Elem())
		}
	}
	return m
}

func (m *MockGormDB) Preload(query string, args ...interface{}) GormDBInterface {
	return m
}

func (m *MockGormDB) Select(query interface{}, args ...interface{}) GormDBInterface {
	return m
}

func (m *MockGormDB) Group(name string) GormDBInterface {
	return m
}

func (m *MockGormDB) Scan(dest interface{}) GormDBInterface {
	if m.ScanValue != nil {
		srcVal := reflect.ValueOf(m.ScanValue)
		dstVal := reflect.ValueOf(dest)

		if srcVal.Kind() == reflect.Slice && dstVal.Kind() == reflect.Ptr && dstVal.Elem().Kind() == reflect.Slice {
			dstVal.Elem().Set(srcVal)
		} else if srcVal.Kind() == reflect.Ptr && dstVal.Kind() == reflect.Ptr && srcVal.Type() == dstVal.Type() {
			dstVal.Elem().Set(srcVal.Elem())
		}
	}
	return m
}

func (m *MockGormDB) Error() error {
	return m.ErrorValue
}

// MockAdminCompanyRepository is a mock implementation of AdminCompanyRepository for testing.
type MockAdminCompanyRepository struct {
	CreateAdminCompanyFunc          func(adminCompany *models.AdminCompaniesTable) error
	GetAdminCompanyByCompanyIDFunc  func(companyID int) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByEmployeeIDFunc func(employeeID int) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByEmailFunc      func(email string) (*models.AdminCompaniesTable, error)
	GetAdminCompanyByIDFunc         func(id int) (*models.AdminCompaniesTable, error)
	UpdateAdminCompanyFunc          func(adminCompany *models.AdminCompaniesTable) error
	GetAdminCompanyByConfirmationTokenFunc func(token string) (*models.AdminCompaniesTable, error)
	ChangeAdminPasswordFunc         func(adminID uint, newPasswordHash string) error
}

func (m *MockAdminCompanyRepository) CreateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	if m.CreateAdminCompanyFunc != nil {
		return m.CreateAdminCompanyFunc(adminCompany)
	}
	return nil
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByCompanyID(companyID int) (*models.AdminCompaniesTable, error) {
	if m.GetAdminCompanyByCompanyIDFunc != nil {
		return m.GetAdminCompanyByCompanyIDFunc(companyID)
	}
	return nil, nil
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByEmployeeID(employeeID int) (*models.AdminCompaniesTable, error) {
	if m.GetAdminCompanyByEmployeeIDFunc != nil {
		return m.GetAdminCompanyByEmployeeIDFunc(employeeID)
	}
	return nil, nil
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByEmail(email string) (*models.AdminCompaniesTable, error) {
	if m.GetAdminCompanyByEmailFunc != nil {
		return m.GetAdminCompanyByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByID(id int) (*models.AdminCompaniesTable, error) {
	if m.GetAdminCompanyByIDFunc != nil {
		return m.GetAdminCompanyByIDFunc(id)
	}
	return nil, nil
}

func (m *MockAdminCompanyRepository) UpdateAdminCompany(adminCompany *models.AdminCompaniesTable) error {
	if m.UpdateAdminCompanyFunc != nil {
		return m.UpdateAdminCompanyFunc(adminCompany)
	}
	return nil
}

func (m *MockAdminCompanyRepository) GetAdminCompanyByConfirmationToken(token string) (*models.AdminCompaniesTable, error) {
	if m.GetAdminCompanyByConfirmationTokenFunc != nil {
		return m.GetAdminCompanyByConfirmationTokenFunc(token)
	}
	return nil, nil
}

func (m *MockAdminCompanyRepository) ChangeAdminPassword(adminID uint, newPasswordHash string) error {
	if m.ChangeAdminPasswordFunc != nil {
		return m.ChangeAdminPasswordFunc(adminID, newPasswordHash)
	}
	return nil
}

var _ repository.AdminCompanyRepository = &MockAdminCompanyRepository{}

// MockAttendanceLocationRepository is a mock implementation of AttendanceLocationRepository for testing.
type MockAttendanceLocationRepository struct {
	mock.Mock
}

func (m *MockAttendanceLocationRepository) CreateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	args := m.Called(location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AttendanceLocation), args.Error(1)
}

func (m *MockAttendanceLocationRepository) GetAttendanceLocationsByCompanyID(companyID uint) ([]models.AttendanceLocation, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AttendanceLocation), args.Error(1)
}

func (m *MockAttendanceLocationRepository) GetAttendanceLocationByID(locationID uint) (*models.AttendanceLocation, error) {
	args := m.Called(locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AttendanceLocation), args.Error(1)
}

func (m *MockAttendanceLocationRepository) GetLocationsByIDs(ids []uint) ([]models.AttendanceLocation, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Ensure the returned slice contains models.AttendanceLocation with ID set correctly
	locations := args.Get(0).([]models.AttendanceLocation)
	for i := range locations {
		locations[i].ID = ids[i] // Assuming IDs match order for simplicity in mock
	}
	return locations, args.Error(1)
}

func (m *MockAttendanceLocationRepository) UpdateAttendanceLocation(location *models.AttendanceLocation) (*models.AttendanceLocation, error) {
	args := m.Called(location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AttendanceLocation), args.Error(1)
}

func (m *MockAttendanceLocationRepository) DeleteAttendanceLocation(locationID uint) error {
	args := m.Called(locationID)
	return args.Error(0)
}

func (m *MockAttendanceLocationRepository) CountAttendanceLocationsByCompanyID(companyID uint) (int64, error) {
	args := m.Called(companyID)
	return args.Get(0).(int64), args.Error(1)
}

var _ repository.AttendanceLocationRepository = &MockAttendanceLocationRepository{}

// MockAttendanceRepository is a mock implementation of AttendanceRepository for testing.
type MockAttendanceRepository struct {
	CreateAttendanceFunc                    func(attendance *models.AttendancesTable) error
	UpdateAttendanceFunc                    func(attendance *models.AttendancesTable) error
	GetLatestAttendanceByEmployeeIDFunc     func(employeeID int) (*models.AttendancesTable, error)
	GetLatestAttendanceForDateFunc          func(employeeID int, date time.Time) (*models.AttendancesTable, error)
	GetLatestOvertimeAttendanceByEmployeeIDFunc func(employeeID int) (*models.AttendancesTable, error)
	GetPresentEmployeesCountTodayFunc       func(companyID int) (int64, error)
	GetAbsentEmployeesCountTodayFunc        func(companyID int) (int64, error)
	GetAttendancesByCompanyIDFunc           func(companyID int) ([]models.AttendancesTable, error)
	GetRecentAttendancesByCompanyIDFunc     func(companyID int, limit int) ([]models.AttendancesTable, error)
	GetRecentOvertimeAttendancesByCompanyIDFunc func(companyID int, limit int) ([]models.AttendancesTable, error)
	GetEmployeeAttendancesFunc              func(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error)
	GetCompanyAttendancesFilteredFunc       func(companyID int, startDate, endDate *time.Time, attendanceType string) ([]models.AttendancesTable, error)
	HasAttendanceForDateFunc                func(employeeID int, date time.Time) (bool, error)
	HasAttendanceForDateRangeFunc           func(employeeID int, startDate, endDate *time.Time) (bool, error)
	GetCompanyOvertimeAttendancesFilteredFunc func(companyID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error)
	GetTodayAttendanceByEmployeeIDFunc      func(employeeID int) (*models.AttendancesTable, error)
	GetRecentAttendancesByEmployeeIDFunc    func(employeeID int, limit int) ([]models.AttendancesTable, error)
	GetAttendancesPaginatedFunc             func(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error)
	GetOvertimeAttendancesPaginatedFunc     func(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error)
	GetUnaccountedEmployeesPaginatedFunc    func(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.EmployeesTable, int64, error)
	GetUnaccountedEmployeesFilteredFunc     func(companyID int, startDate, endDate *time.Time, search string) ([]models.EmployeesTable, error)
	GetOvertimeAttendancesFilteredFunc      func(companyID int, startDate, endDate *time.Time, search string) ([]models.AttendancesTable, error)
}

func (m *MockAttendanceRepository) CreateAttendance(attendance *models.AttendancesTable) error {
	if m.CreateAttendanceFunc != nil {
		return m.CreateAttendanceFunc(attendance)
	}
	return nil
}

func (m *MockAttendanceRepository) UpdateAttendance(attendance *models.AttendancesTable) error {
	if m.UpdateAttendanceFunc != nil {
		return m.UpdateAttendanceFunc(attendance)
	}
	return nil
}

func (m *MockAttendanceRepository) GetLatestAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	if m.GetLatestAttendanceByEmployeeIDFunc != nil {
		return m.GetLatestAttendanceByEmployeeIDFunc(employeeID)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetLatestAttendanceForDate(employeeID int, date time.Time) (*models.AttendancesTable, error) {
	if m.GetLatestAttendanceForDateFunc != nil {
		return m.GetLatestAttendanceForDateFunc(employeeID, date)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetLatestOvertimeAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	if m.GetLatestOvertimeAttendanceByEmployeeIDFunc != nil {
		return m.GetLatestOvertimeAttendanceByEmployeeIDFunc(employeeID)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetPresentEmployeesCountToday(companyID int) (int64, error) {
	if m.GetPresentEmployeesCountTodayFunc != nil {
		return m.GetPresentEmployeesCountTodayFunc(companyID)
	}
	return 0, nil
}

func (m *MockAttendanceRepository) GetAbsentEmployeesCountToday(companyID int) (int64, error) {
	if m.GetAbsentEmployeesCountTodayFunc != nil {
		return m.GetAbsentEmployeesCountTodayFunc(companyID)
	}
	return 0, nil
}

func (m *MockAttendanceRepository) GetAttendancesByCompanyID(companyID int) ([]models.AttendancesTable, error) {
	if m.GetAttendancesByCompanyIDFunc != nil {
		return m.GetAttendancesByCompanyIDFunc(companyID)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetRecentAttendancesByCompanyID(companyID int, limit int) ([]models.AttendancesTable, error) {
	if m.GetRecentAttendancesByCompanyIDFunc != nil {
		return m.GetRecentAttendancesByCompanyIDFunc(companyID, limit)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetRecentOvertimeAttendancesByCompanyID(companyID int, limit int) ([]models.AttendancesTable, error) {
	if m.GetRecentOvertimeAttendancesByCompanyIDFunc != nil {
		return m.GetRecentOvertimeAttendancesByCompanyIDFunc(companyID, limit)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetEmployeeAttendances(employeeID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	if m.GetEmployeeAttendancesFunc != nil {
		return m.GetEmployeeAttendancesFunc(employeeID, startDate, endDate)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetCompanyAttendancesFiltered(companyID int, startDate, endDate *time.Time, attendanceType string) ([]models.AttendancesTable, error) {
	if m.GetCompanyAttendancesFilteredFunc != nil {
		return m.GetCompanyAttendancesFilteredFunc(companyID, startDate, endDate, attendanceType)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) HasAttendanceForDate(employeeID int, date time.Time) (bool, error) {
	if m.HasAttendanceForDateFunc != nil {
		return m.HasAttendanceForDateFunc(employeeID, date)
	}
	return false, nil
}

func (m *MockAttendanceRepository) HasAttendanceForDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	if m.HasAttendanceForDateRangeFunc != nil {
		return m.HasAttendanceForDateRangeFunc(employeeID, startDate, endDate)
	}
	return false, nil
}

func (m *MockAttendanceRepository) GetCompanyOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time) ([]models.AttendancesTable, error) {
	if m.GetCompanyOvertimeAttendancesFilteredFunc != nil {
		return m.GetCompanyOvertimeAttendancesFilteredFunc(companyID, startDate, endDate)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetTodayAttendanceByEmployeeID(employeeID int) (*models.AttendancesTable, error) {
	if m.GetTodayAttendanceByEmployeeIDFunc != nil {
		return m.GetTodayAttendanceByEmployeeIDFunc(employeeID)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetRecentAttendancesByEmployeeID(employeeID int, limit int) ([]models.AttendancesTable, error) {
	if m.GetRecentAttendancesByEmployeeIDFunc != nil {
		return m.GetRecentAttendancesByEmployeeIDFunc(employeeID, limit)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error) {
	if m.GetAttendancesPaginatedFunc != nil {
		return m.GetAttendancesPaginatedFunc(companyID, startDate, endDate, search, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockAttendanceRepository) GetOvertimeAttendancesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.AttendancesTable, int64, error) {
	if m.GetOvertimeAttendancesPaginatedFunc != nil {
		return m.GetOvertimeAttendancesPaginatedFunc(companyID, startDate, endDate, search, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockAttendanceRepository) GetUnaccountedEmployeesPaginated(companyID int, startDate, endDate *time.Time, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	if m.GetUnaccountedEmployeesPaginatedFunc != nil {
		return m.GetUnaccountedEmployeesPaginatedFunc(companyID, startDate, endDate, search, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockAttendanceRepository) GetUnaccountedEmployeesFiltered(companyID int, startDate, endDate *time.Time, search string) ([]models.EmployeesTable, error) {
	if m.GetUnaccountedEmployeesFilteredFunc != nil {
		return m.GetUnaccountedEmployeesFilteredFunc(companyID, startDate, endDate, search)
	}
	return nil, nil
}

func (m *MockAttendanceRepository) GetOvertimeAttendancesFiltered(companyID int, startDate, endDate *time.Time, search string) ([]models.AttendancesTable, error) {
	if m.GetOvertimeAttendancesFilteredFunc != nil {
		return m.GetOvertimeAttendancesFilteredFunc(companyID, startDate, endDate, search)
	}
	return nil, nil
}

var _ repository.AttendanceRepository = &MockAttendanceRepository{}

// MockBroadcastRepository is a mock implementation of BroadcastRepository for testing.
type MockBroadcastRepository struct {
	CreateBroadcastFunc      func(message *models.BroadcastMessage) error
	GetBroadcastsForEmployeeFunc func(companyID, employeeID uint) ([]models.BroadcastMessage, error)
	MarkBroadcastAsReadFunc  func(employeeID, messageID uint) error
}

func (m *MockBroadcastRepository) CreateBroadcast(message *models.BroadcastMessage) error {
	if m.CreateBroadcastFunc != nil {
		return m.CreateBroadcastFunc(message)
	}
	return nil
}

func (m *MockBroadcastRepository) GetBroadcastsForEmployee(companyID, employeeID uint) ([]models.BroadcastMessage, error) {
	if m.GetBroadcastsForEmployeeFunc != nil {
		return m.GetBroadcastsForEmployeeFunc(companyID, employeeID)
	}
	return nil, nil
}

func (m *MockBroadcastRepository) MarkBroadcastAsRead(employeeID, messageID uint) error {
	if m.MarkBroadcastAsReadFunc != nil {
		return m.MarkBroadcastAsReadFunc(employeeID, messageID)
	}
	return nil
}

var _ repository.BroadcastRepository = &MockBroadcastRepository{}

// MockCompanyRepository is a mock implementation of CompanyRepository for testing.
type MockCompanyRepository struct {
	CreateCompanyFunc                 func(company *models.CompaniesTable) error
	GetCompanyByIDFunc                func(id int) (*models.CompaniesTable, error)
	GetCompanyWithSubscriptionDetailsFunc func(id int) (*models.CompaniesTable, error)
	UpdateCompanyFunc                 func(company *models.CompaniesTable) error
	GetAllActiveCompaniesFunc         func() ([]models.CompaniesTable, error)
	CreateCompanyWithAdminAndShiftFunc func(company *models.CompaniesTable, admin *models.AdminCompaniesTable, shift *models.ShiftsTable) error
	DeleteCompanyFunc                 func(id int) error
	GetTotalEmployeesByCompanyIDFunc  func(companyID int) (int64, error)
	GetAllCompaniesFunc               func() ([]models.CompaniesTable, error)
}

func (m *MockCompanyRepository) CreateCompany(company *models.CompaniesTable) error {
	if m.CreateCompanyFunc != nil {
		return m.CreateCompanyFunc(company)
	}
	return nil
}

func (m *MockCompanyRepository) GetCompanyByID(id int) (*models.CompaniesTable, error) {
	if m.GetCompanyByIDFunc != nil {
		return m.GetCompanyByIDFunc(id)
	}
	return nil, nil
}

func (m *MockCompanyRepository) GetCompanyWithSubscriptionDetails(id int) (*models.CompaniesTable, error) {
	if m.GetCompanyWithSubscriptionDetailsFunc != nil {
		return m.GetCompanyWithSubscriptionDetailsFunc(id)
	}
	return nil, nil
}

func (m *MockCompanyRepository) UpdateCompany(company *models.CompaniesTable) error {
	if m.UpdateCompanyFunc != nil {
		return m.UpdateCompanyFunc(company)
	}
	return nil
}

func (m *MockCompanyRepository) GetAllActiveCompanies() ([]models.CompaniesTable, error) {
	if m.GetAllActiveCompaniesFunc != nil {
		return m.GetAllActiveCompaniesFunc()
	}
	return nil, nil
}

func (m *MockCompanyRepository) CreateCompanyWithAdminAndShift(company *models.CompaniesTable, admin *models.AdminCompaniesTable, shift *models.ShiftsTable) error {
	if m.CreateCompanyWithAdminAndShiftFunc != nil {
		return m.CreateCompanyWithAdminAndShiftFunc(company, admin, shift)
	}
	return nil
}

func (m *MockCompanyRepository) DeleteCompany(id int) error {
	if m.DeleteCompanyFunc != nil {
		return m.DeleteCompanyFunc(id)
	}
	return nil
}

func (m *MockCompanyRepository) GetTotalEmployeesByCompanyID(companyID int) (int64, error) {
	if m.GetTotalEmployeesByCompanyIDFunc != nil {
		return m.GetTotalEmployeesByCompanyIDFunc(companyID)
	}
	return 0, nil
}

func (m *MockCompanyRepository) GetAllCompanies() ([]models.CompaniesTable, error) {
	if m.GetAllCompaniesFunc != nil {
		return m.GetAllCompaniesFunc()
	}
	return nil, nil
}

var _ repository.CompanyRepository = &MockCompanyRepository{}

// MockCustomOfferRepository is a mock implementation of CustomOfferRepository for testing.
type MockCustomOfferRepository struct {
	CreateCustomOfferFunc   func(offer *models.CustomOffer) error
	GetCustomOfferByTokenFunc func(token string) (*models.CustomOffer, error)
	GetCustomOfferByIDFunc  func(id uint) (*models.CustomOffer, error)
	UpdateCustomOfferFunc   func(offer *models.CustomOffer) error
	MarkCustomOfferAsUsedFunc func(token string) error
}

func (m *MockCustomOfferRepository) CreateCustomOffer(offer *models.CustomOffer) error {
	if m.CreateCustomOfferFunc != nil {
		return m.CreateCustomOfferFunc(offer)
	}
	return nil
}

func (m *MockCustomOfferRepository) GetCustomOfferByToken(token string) (*models.CustomOffer, error) {
	if m.GetCustomOfferByTokenFunc != nil {
		return m.GetCustomOfferByTokenFunc(token)
	}
	return nil, nil
}

func (m *MockCustomOfferRepository) GetCustomOfferByID(id uint) (*models.CustomOffer, error) {
	if m.GetCustomOfferByIDFunc != nil {
		return m.GetCustomOfferByIDFunc(id)
	}
	return nil, nil
}

func (m *MockCustomOfferRepository) UpdateCustomOffer(offer *models.CustomOffer) error {
	if m.UpdateCustomOfferFunc != nil {
		return m.UpdateCustomOfferFunc(offer)
	}
	return nil
}

func (m *MockCustomOfferRepository) MarkCustomOfferAsUsed(token string) error {
	if m.MarkCustomOfferAsUsedFunc != nil {
		return m.MarkCustomOfferAsUsedFunc(token)
	}
	return nil
}

var _ repository.CustomOfferRepository = &MockCustomOfferRepository{}

// MockCustomPackageRequestRepository is a mock implementation of CustomPackageRequestRepository for testing.
type MockCustomPackageRequestRepository struct {
	CreateCustomPackageRequestFunc      func(req *models.CustomPackageRequest) error
	GetCustomPackageRequestsPaginatedFunc func(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error)
	GetCustomPackageRequestByIDFunc     func(id uint) (*models.CustomPackageRequest, error)
	UpdateCustomPackageRequestFunc      func(req *models.CustomPackageRequest) error
}

func (m *MockCustomPackageRequestRepository) CreateCustomPackageRequest(req *models.CustomPackageRequest) error {
	if m.CreateCustomPackageRequestFunc != nil {
		return m.CreateCustomPackageRequestFunc(req)
	}
	return nil
}

func (m *MockCustomPackageRequestRepository) GetCustomPackageRequestsPaginated(page, pageSize int, search string) ([]models.CustomPackageRequest, int64, error) {
	if m.GetCustomPackageRequestsPaginatedFunc != nil {
		return m.GetCustomPackageRequestsPaginatedFunc(page, pageSize, search)
	}
	return nil, 0, nil
}

func (m *MockCustomPackageRequestRepository) GetCustomPackageRequestByID(id uint) (*models.CustomPackageRequest, error) {
	if m.GetCustomPackageRequestByIDFunc != nil {
		return m.GetCustomPackageRequestByIDFunc(id)
	}
	return nil, nil
}

func (m *MockCustomPackageRequestRepository) UpdateCustomPackageRequest(req *models.CustomPackageRequest) error {
	if m.UpdateCustomPackageRequestFunc != nil {
		return m.UpdateCustomPackageRequestFunc(req)
	}
	return nil
}

var _ repository.CustomPackageRequestRepository = &MockCustomPackageRequestRepository{}

// MockDivisionRepository is a mock implementation of DivisionRepository for testing.
type MockDivisionRepository struct {
	mock.Mock
}

// We implement the interface methods on the mock struct.
func (m *MockDivisionRepository) CreateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	args := m.Called(division)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) IsDivisionNameTaken(name string, companyID uint, currentDivisionID uint) (bool, error) {
	args := m.Called(name, companyID, currentDivisionID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDivisionRepository) GetDivisionsByCompanyID(companyID uint) ([]models.DivisionTable, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) GetDivisionByID(divisionID uint) (*models.DivisionTable, error) {
	args := m.Called(divisionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) UpdateDivision(division *models.DivisionTable) (*models.DivisionTable, error) {
	args := m.Called(division)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DivisionTable), args.Error(1)
}

func (m *MockDivisionRepository) DeleteDivision(divisionID uint) error {
	args := m.Called(divisionID)
	return args.Error(0)
}

var _ repository.DivisionRepository = &MockDivisionRepository{}



// MockEmployeeRepository is a mock implementation of EmployeeRepository for testing.
type MockEmployeeRepository struct {
	CreateEmployeeFunc                   func(employee *models.EmployeesTable) error
	GetEmployeeByIDFunc                  func(id int) (*models.EmployeesTable, error)
	GetEmployeesByCompanyIDFunc          func(companyID int) ([]models.EmployeesTable, error)
	GetEmployeeByEmailFunc               func(email string) (*models.EmployeesTable, error)
	UpdateEmployeeFunc                   func(employee *models.EmployeesTable) error
	DeleteEmployeeFunc                   func(id int) error
	SearchEmployeesFunc                  func(companyID int, name string) ([]models.EmployeesTable, error)
	GetTotalEmployeesByCompanyIDFunc     func(companyID int) (int64, error)
	GetEmployeesWithFaceImagesFunc       func(companyID int) ([]models.EmployeesTable, error)
	GetOnLeaveEmployeesCountTodayFunc    func(companyID int) (int64, error)
	SetEmployeePasswordSetFunc           func(employeeID uint, isSet bool) error
	UpdateEmployeePasswordFunc           func(employee *models.EmployeesTable, newPassword string) error
	GetPendingEmployeesFunc              func(companyID int) ([]models.EmployeesTable, error)
	GetEmployeeByEmailOrIDNumberFunc     func(email, employeeIDNumber string) (*models.EmployeesTable, error)
	GetEmployeesByCompanyIDPaginatedFunc func(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error)
	GetPendingEmployeesByCompanyIDPaginatedFunc func(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error)
	UpdateEmployeeFieldsFunc             func(employee *models.EmployeesTable, updates map[string]interface{}) error
	GetActiveEmployeesByCompanyIDFunc    func(companyID int) ([]models.EmployeesTable, error)
	GetAbsentEmployeesCountTodayFunc     func(companyID int, presentToday int64) (int64, error)
}

func (m *MockEmployeeRepository) CreateEmployee(employee *models.EmployeesTable) error {
	if m.CreateEmployeeFunc != nil {
		return m.CreateEmployeeFunc(employee)
	}
	return nil
}

func (m *MockEmployeeRepository) GetEmployeeByID(id int) (*models.EmployeesTable, error) {
	if m.GetEmployeeByIDFunc != nil {
		return m.GetEmployeeByIDFunc(id)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error) {
	if m.GetEmployeesByCompanyIDFunc != nil {
		return m.GetEmployeesByCompanyIDFunc(companyID)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetEmployeeByEmail(email string) (*models.EmployeesTable, error) {
	if m.GetEmployeeByEmailFunc != nil {
		return m.GetEmployeeByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) UpdateEmployee(employee *models.EmployeesTable) error {
	if m.UpdateEmployeeFunc != nil {
		return m.UpdateEmployeeFunc(employee)
	}
	return nil
}

func (m *MockEmployeeRepository) DeleteEmployee(id int) error {
	if m.DeleteEmployeeFunc != nil {
		return m.DeleteEmployeeFunc(id)
	}
	return nil
}

func (m *MockEmployeeRepository) SearchEmployees(companyID int, name string) ([]models.EmployeesTable, error) {
	if m.SearchEmployeesFunc != nil {
		return m.SearchEmployeesFunc(companyID, name)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetTotalEmployeesByCompanyID(companyID int) (int64, error) {
	if m.GetTotalEmployeesByCompanyIDFunc != nil {
		return m.GetTotalEmployeesByCompanyIDFunc(companyID)
	}
	return 0, nil
}

func (m *MockEmployeeRepository) GetEmployeesWithFaceImages(companyID int) ([]models.EmployeesTable, error) {
	if m.GetEmployeesWithFaceImagesFunc != nil {
		return m.GetEmployeesWithFaceImagesFunc(companyID)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetOnLeaveEmployeesCountToday(companyID int) (int64, error) {
	if m.GetOnLeaveEmployeesCountTodayFunc != nil {
		return m.GetOnLeaveEmployeesCountTodayFunc(companyID)
	}
	return 0, nil
}

func (m *MockEmployeeRepository) SetEmployeePasswordSet(employeeID uint, isSet bool) error {
	if m.SetEmployeePasswordSetFunc != nil {
		return m.SetEmployeePasswordSetFunc(employeeID, isSet)
	}
	return nil
}

func (m *MockEmployeeRepository) UpdateEmployeePassword(employee *models.EmployeesTable, newPassword string) error {
	if m.UpdateEmployeePasswordFunc != nil {
		return m.UpdateEmployeePasswordFunc(employee, newPassword)
	}
	return nil
}

func (m *MockEmployeeRepository) GetPendingEmployees(companyID int) ([]models.EmployeesTable, error) {
	if m.GetPendingEmployeesFunc != nil {
		return m.GetPendingEmployeesFunc(companyID)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetEmployeeByEmailOrIDNumber(email, employeeIDNumber string) (*models.EmployeesTable, error) {
	if m.GetEmployeeByEmailOrIDNumberFunc != nil {
		return m.GetEmployeeByEmailOrIDNumberFunc(email, employeeIDNumber)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	if m.GetEmployeesByCompanyIDPaginatedFunc != nil {
		return m.GetEmployeesByCompanyIDPaginatedFunc(companyID, search, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockEmployeeRepository) GetPendingEmployeesByCompanyIDPaginated(companyID int, search string, page, pageSize int) ([]models.EmployeesTable, int64, error) {
	if m.GetPendingEmployeesByCompanyIDPaginatedFunc != nil {
		return m.GetPendingEmployeesByCompanyIDPaginatedFunc(companyID, search, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockEmployeeRepository) UpdateEmployeeFields(employee *models.EmployeesTable, updates map[string]interface{}) error {
	if m.UpdateEmployeeFieldsFunc != nil {
		return m.UpdateEmployeeFieldsFunc(employee, updates)
	}
	return nil
}

func (m *MockEmployeeRepository) GetActiveEmployeesByCompanyID(companyID int) ([]models.EmployeesTable, error) {
	if m.GetActiveEmployeesByCompanyIDFunc != nil {
		return m.GetActiveEmployeesByCompanyIDFunc(companyID)
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetAbsentEmployeesCountToday(companyID int, presentToday int64) (int64, error) {
	if m.GetAbsentEmployeesCountTodayFunc != nil {
		return m.GetAbsentEmployeesCountTodayFunc(companyID, presentToday)
	}
	return 0, nil
}

var _ repository.EmployeeRepository = &MockEmployeeRepository{}

// MockFaceImageRepository is a mock implementation of FaceImageRepository for testing.
type MockFaceImageRepository struct {
	CreateFaceImageFunc       func(faceImage *models.FaceImagesTable) error
	GetFaceImagesByEmployeeIDFunc func(employeeID int) ([]models.FaceImagesTable, error)
	GetFaceImageByIDFunc      func(id int) (*models.FaceImagesTable, error)
	DeleteFaceImageFunc       func(id int) error
}

func (m *MockFaceImageRepository) CreateFaceImage(faceImage *models.FaceImagesTable) error {
	if m.CreateFaceImageFunc != nil {
		return m.CreateFaceImageFunc(faceImage)
	}
	return nil
}

func (m *MockFaceImageRepository) GetFaceImagesByEmployeeID(employeeID int) ([]models.FaceImagesTable, error) {
	if m.GetFaceImagesByEmployeeIDFunc != nil {
		return m.GetFaceImagesByEmployeeIDFunc(employeeID)
	}
	return nil, nil
}

func (m *MockFaceImageRepository) GetFaceImageByID(id int) (*models.FaceImagesTable, error) {
	if m.GetFaceImageByIDFunc != nil {
		return m.GetFaceImageByIDFunc(id)
	}
	return nil, nil
}

func (m *MockFaceImageRepository) DeleteFaceImage(id int) error {
	if m.DeleteFaceImageFunc != nil {
		return m.DeleteFaceImageFunc(id)
	}
	return nil
}

var _ repository.FaceImageRepository = &MockFaceImageRepository{}

// MockInvoiceRepository is a mock implementation of InvoiceRepository for testing.
type MockInvoiceRepository struct {
	CreateInvoiceFunc        func(invoice *models.InvoiceTable) error
	GetInvoiceByOrderIDFunc  func(orderID string) (*models.InvoiceTable, error)
	GetInvoicesByCompanyIDFunc func(companyID uint) ([]models.InvoiceTable, error)
	UpdateInvoiceFunc        func(invoice *models.InvoiceTable) error
}

func (m *MockInvoiceRepository) CreateInvoice(invoice *models.InvoiceTable) error {
	if m.CreateInvoiceFunc != nil {
		return m.CreateInvoiceFunc(invoice)
	}
	return nil
}

func (m *MockInvoiceRepository) GetInvoiceByOrderID(orderID string) (*models.InvoiceTable, error) {
	if m.GetInvoiceByOrderIDFunc != nil {
		return m.GetInvoiceByOrderIDFunc(orderID)
	}
	return nil, nil
}

func (m *MockInvoiceRepository) GetInvoicesByCompanyID(companyID uint) ([]models.InvoiceTable, error) {
	if m.GetInvoicesByCompanyIDFunc != nil {
		return m.GetInvoicesByCompanyIDFunc(companyID)
	}
	return nil, nil
}

func (m *MockInvoiceRepository) UpdateInvoice(invoice *models.InvoiceTable) error {
	if m.UpdateInvoiceFunc != nil {
		return m.UpdateInvoiceFunc(invoice)
	}
	return nil
}

var _ repository.InvoiceRepository = &MockInvoiceRepository{}

// MockLeaveRequestRepository is a mock implementation of LeaveRequestRepository for testing.
type MockLeaveRequestRepository struct {
	CreateLeaveRequestFunc              func(leaveRequest *models.LeaveRequest) error
	GetLeaveRequestByIDFunc             func(id uint) (*models.LeaveRequest, error)
	GetAllLeaveRequestsFunc             func() ([]models.LeaveRequest, error)
	GetLeaveRequestsByEmployeeIDFunc    func(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error)
	GetCompanyLeaveRequestsFilteredFunc func(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error)
	UpdateLeaveRequestFunc              func(leaveRequest *models.LeaveRequest) error
	GetRecentLeaveRequestsByCompanyIDFunc func(companyID int, limit int) ([]models.LeaveRequest, error)
	IsEmployeeOnApprovedLeaveFunc       func(employeeID int, date time.Time) (*models.LeaveRequest, error)
	IsEmployeeOnApprovedLeaveDateRangeFunc func(employeeID int, startDate, endDate *time.Time) (bool, error)
	GetPendingLeaveRequestsByEmployeeIDFunc func(employeeID int) ([]models.LeaveRequest, error)
	GetCompanyLeaveRequestsPaginatedFunc func(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error)
	GetOnLeaveEmployeesCountTodayFunc   func(companyID int) (int64, error)
}

func (m *MockLeaveRequestRepository) CreateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	if m.CreateLeaveRequestFunc != nil {
		return m.CreateLeaveRequestFunc(leaveRequest)
	}
	return nil
}

func (m *MockLeaveRequestRepository) GetLeaveRequestByID(id uint) (*models.LeaveRequest, error) {
	if m.GetLeaveRequestByIDFunc != nil {
		return m.GetLeaveRequestByIDFunc(id)
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) GetAllLeaveRequests() ([]models.LeaveRequest, error) {
	if m.GetAllLeaveRequestsFunc != nil {
		return m.GetAllLeaveRequestsFunc()
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) GetLeaveRequestsByEmployeeID(employeeID uint, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	if m.GetLeaveRequestsByEmployeeIDFunc != nil {
		return m.GetLeaveRequestsByEmployeeIDFunc(employeeID, startDate, endDate)
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) GetCompanyLeaveRequestsFiltered(companyID int, status, search string, startDate, endDate *time.Time) ([]models.LeaveRequest, error) {
	if m.GetCompanyLeaveRequestsFilteredFunc != nil {
		return m.GetCompanyLeaveRequestsFilteredFunc(companyID, status, search, startDate, endDate)
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) UpdateLeaveRequest(leaveRequest *models.LeaveRequest) error {
	if m.UpdateLeaveRequestFunc != nil {
		return m.UpdateLeaveRequestFunc(leaveRequest)
	}
	return nil
}

func (m *MockLeaveRequestRepository) GetRecentLeaveRequestsByCompanyID(companyID int, limit int) ([]models.LeaveRequest, error) {
	if m.GetRecentLeaveRequestsByCompanyIDFunc != nil {
		return m.GetRecentLeaveRequestsByCompanyIDFunc(companyID, limit)
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) IsEmployeeOnApprovedLeave(employeeID int, date time.Time) (*models.LeaveRequest, error) {
	if m.IsEmployeeOnApprovedLeaveFunc != nil {
		return m.IsEmployeeOnApprovedLeaveFunc(employeeID, date)
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) IsEmployeeOnApprovedLeaveDateRange(employeeID int, startDate, endDate *time.Time) (bool, error) {
	if m.IsEmployeeOnApprovedLeaveDateRangeFunc != nil {
		return m.IsEmployeeOnApprovedLeaveDateRangeFunc(employeeID, startDate, endDate)
	}
	return false, nil
}

func (m *MockLeaveRequestRepository) GetPendingLeaveRequestsByEmployeeID(employeeID int) ([]models.LeaveRequest, error) {
	if m.GetPendingLeaveRequestsByEmployeeIDFunc != nil {
		return m.GetPendingLeaveRequestsByEmployeeIDFunc(employeeID)
	}
	return nil, nil
}

func (m *MockLeaveRequestRepository) GetCompanyLeaveRequestsPaginated(companyID int, status, search string, startDate, endDate *time.Time, page, pageSize int) ([]models.LeaveRequest, int64, error) {
	if m.GetCompanyLeaveRequestsPaginatedFunc != nil {
		return m.GetCompanyLeaveRequestsPaginatedFunc(companyID, status, search, startDate, endDate, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockLeaveRequestRepository) GetOnLeaveEmployeesCountToday(companyID int) (int64, error) {
	if m.GetOnLeaveEmployeesCountTodayFunc != nil {
		return m.GetOnLeaveEmployeesCountTodayFunc(companyID)
	}
	return 0, nil
}

var _ repository.LeaveRequestRepository = &MockLeaveRequestRepository{}

// MockPasswordResetRepository is a mock implementation of PasswordResetRepository for testing.
type MockPasswordResetRepository struct {
	CreatePasswordResetTokenFunc      func(token *models.PasswordResetTokenTable) error
	GetPasswordResetTokenFunc         func(tokenString string) (*models.PasswordResetTokenTable, error)
	MarkPasswordResetTokenAsUsedFunc  func(token *models.PasswordResetTokenTable) error
	InvalidatePasswordResetTokensByUserIDFunc func(userID uint, tokenType string) error
}

func (m *MockPasswordResetRepository) CreatePasswordResetToken(token *models.PasswordResetTokenTable) error {
	if m.CreatePasswordResetTokenFunc != nil {
		return m.CreatePasswordResetTokenFunc(token)
	}
	return nil
}

func (m *MockPasswordResetRepository) GetPasswordResetToken(tokenString string) (*models.PasswordResetTokenTable, error) {
	if m.GetPasswordResetTokenFunc != nil {
		return m.GetPasswordResetTokenFunc(tokenString)
	}
	return nil, nil
}

func (m *MockPasswordResetRepository) MarkPasswordResetTokenAsUsed(token *models.PasswordResetTokenTable) error {
	if m.MarkPasswordResetTokenAsUsedFunc != nil {
		return m.MarkPasswordResetTokenAsUsedFunc(token)
	}
	return nil
}

func (m *MockPasswordResetRepository) InvalidatePasswordResetTokensByUserID(userID uint, tokenType string) error {
	if m.InvalidatePasswordResetTokensByUserIDFunc != nil {
		return m.InvalidatePasswordResetTokensByUserIDFunc(userID, tokenType)
	}
	return nil
}

var _ repository.PasswordResetRepository = &MockPasswordResetRepository{}

// MockShiftRepository is a mock implementation of ShiftRepository for testing.
type MockShiftRepository struct {
	mock.Mock
}

func (m *MockShiftRepository) CreateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	args := m.Called(shift)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShiftsTable), args.Error(1)
}

func (m *MockShiftRepository) GetShiftsByCompanyID(companyID int) ([]models.ShiftsTable, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ShiftsTable), args.Error(1)
}

func (m *MockShiftRepository) GetShiftByID(id int) (*models.ShiftsTable, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShiftsTable), args.Error(1)
}

func (m *MockShiftRepository) GetShiftsByIDs(ids []uint) ([]models.ShiftsTable, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Ensure the returned slice contains models.ShiftsTable with ID set correctly
	shifts := args.Get(0).([]models.ShiftsTable)
	for i := range shifts {
		shifts[i].ID = int(ids[i]) // Assuming IDs match order for simplicity in mock
	}
	return shifts, args.Error(1)
}

func (m *MockShiftRepository) UpdateShift(shift *models.ShiftsTable) (*models.ShiftsTable, error) {
	args := m.Called(shift)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShiftsTable), args.Error(1)
}

func (m *MockShiftRepository) DeleteShift(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockShiftRepository) SetDefaultShift(companyID, shiftID int) error {
	args := m.Called(companyID, shiftID)
	return args.Error(0)
}

func (m *MockShiftRepository) GetDefaultShiftByCompanyID(companyID int) (*models.ShiftsTable, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShiftsTable), args.Error(1)
}

var _ repository.ShiftRepository = &MockShiftRepository{}

// MockSubscriptionPackageRepository is a mock implementation of SubscriptionPackageRepository for testing.
type MockSubscriptionPackageRepository struct {
	GetSubscriptionPackagesFunc     func() ([]models.SubscriptionPackageTable, error)
	CreateSubscriptionPackageFunc     func(pkg *models.SubscriptionPackageTable) error
	UpdateSubscriptionPackageFunc     func(pkg *models.SubscriptionPackageTable) error
	UpdateSubscriptionPackageFieldsFunc func(id string, updates map[string]interface{}) error
	DeleteSubscriptionPackageFunc     func(id string) error
	GetSubscriptionPackageByIDFunc    func(id int) (*models.SubscriptionPackageTable, error)
}

func (m *MockSubscriptionPackageRepository) GetSubscriptionPackages() ([]models.SubscriptionPackageTable, error) {
	if m.GetSubscriptionPackagesFunc != nil {
		return m.GetSubscriptionPackagesFunc()
	}
	return nil, nil
}

func (m *MockSubscriptionPackageRepository) CreateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error {
	if m.CreateSubscriptionPackageFunc != nil {
		return m.CreateSubscriptionPackageFunc(pkg)
	}
	return nil
}

func (m *MockSubscriptionPackageRepository) UpdateSubscriptionPackage(pkg *models.SubscriptionPackageTable) error {
	if m.UpdateSubscriptionPackageFunc != nil {
		return m.UpdateSubscriptionPackageFunc(pkg)
	}
	return nil
}

func (m *MockSubscriptionPackageRepository) UpdateSubscriptionPackageFields(id string, updates map[string]interface{}) error {
	if m.UpdateSubscriptionPackageFieldsFunc != nil {
		return m.UpdateSubscriptionPackageFieldsFunc(id, updates)
	}
	return nil
}

func (m *MockSubscriptionPackageRepository) DeleteSubscriptionPackage(id string) error {
	if m.DeleteSubscriptionPackageFunc != nil {
		return m.DeleteSubscriptionPackageFunc(id)
	}
	return nil
}

func (m *MockSubscriptionPackageRepository) GetSubscriptionPackageByID(id int) (*models.SubscriptionPackageTable, error) {
	if m.GetSubscriptionPackageByIDFunc != nil {
		return m.GetSubscriptionPackageByIDFunc(id)
	}
	return nil, nil
}

var _ repository.SubscriptionPackageRepository = &MockSubscriptionPackageRepository{}

// MockSuperAdminRepository is a mock implementation of SuperAdminRepository for testing.
type MockSuperAdminRepository struct {
	CreateSuperAdminFunc  func(superUser *models.SuperAdminTable) error
	GetSuperAdminByIDFunc   func(id int) (*models.SuperAdminTable, error)
	GetSuperAdminByEmailFunc func(email string) (*models.SuperAdminTable, error)
	UpdateSuperAdminFunc  func(id int, superUser *models.SuperAdminTable) (*models.SuperAdminTable, error)
	DeleteSuperAdminFunc  func(id int) error
	GetAllSuperAdminsFunc func() ([]models.SuperAdminTable, error)
	GetTotalCompaniesCountFunc func() (int64, error)
	GetCompaniesCountBySubscriptionStatusFunc func(status string) (int64, error)
	GetExpiredAndTrialExpiredCompaniesCountFunc func() (int64, error)
	GetRecentCompaniesFunc func(limit int) ([]models.CompaniesTable, error)
	GetAllCompaniesWithPreloadFunc func() ([]models.CompaniesTable, error)
	GetPaidInvoicesMonthlyRevenueFunc func(startDate, endDate *time.Time) ([]struct {
		Month        string
		Year         string
		TotalRevenue float64
	}, error)
}

func (m *MockSuperAdminRepository) CreateSuperAdmin(superUser *models.SuperAdminTable) error {
	if m.CreateSuperAdminFunc != nil {
		return m.CreateSuperAdminFunc(superUser)
	}
	return nil
}

func (m *MockSuperAdminRepository) GetSuperAdminByID(id int) (*models.SuperAdminTable, error) {
	if m.GetSuperAdminByIDFunc != nil {
		return m.GetSuperAdminByIDFunc(id)
	}
	return nil, nil
}

func (m *MockSuperAdminRepository) GetSuperAdminByEmail(email string) (*models.SuperAdminTable, error) {
	if m.GetSuperAdminByEmailFunc != nil {
		return m.GetSuperAdminByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockSuperAdminRepository) UpdateSuperAdmin(id int, superUser *models.SuperAdminTable) (*models.SuperAdminTable, error) {
	if m.UpdateSuperAdminFunc != nil {
		return m.UpdateSuperAdminFunc(id, superUser)
	}
	return nil, nil
}

func (m *MockSuperAdminRepository) DeleteSuperAdmin(id int) error {
	if m.DeleteSuperAdminFunc != nil {
		return m.DeleteSuperAdminFunc(id)
	}
	return nil
}

func (m *MockSuperAdminRepository) GetAllSuperAdmins() ([]models.SuperAdminTable, error) {
	if m.GetAllSuperAdminsFunc != nil {
		return m.GetAllSuperAdminsFunc()
	}
	return nil, nil
}

func (m *MockSuperAdminRepository) GetTotalCompaniesCount() (int64, error) {
	if m.GetTotalCompaniesCountFunc != nil {
		return m.GetTotalCompaniesCountFunc()
	}
	return 0, nil
}

func (m *MockSuperAdminRepository) GetCompaniesCountBySubscriptionStatus(status string) (int64, error) {
	if m.GetCompaniesCountBySubscriptionStatusFunc != nil {
		return m.GetCompaniesCountBySubscriptionStatusFunc(status)
	}
	return 0, nil
}

func (m *MockSuperAdminRepository) GetExpiredAndTrialExpiredCompaniesCount() (int64, error) {
	if m.GetExpiredAndTrialExpiredCompaniesCountFunc != nil {
		return m.GetExpiredAndTrialExpiredCompaniesCountFunc()
	}
	return 0, nil
}

func (m *MockSuperAdminRepository) GetRecentCompanies(limit int) ([]models.CompaniesTable, error) {
	if m.GetRecentCompaniesFunc != nil {
		return m.GetRecentCompaniesFunc(limit)
	}
	return nil, nil
}

func (m *MockSuperAdminRepository) GetAllCompaniesWithPreload() ([]models.CompaniesTable, error) {
	if m.GetAllCompaniesWithPreloadFunc != nil {
		return m.GetAllCompaniesWithPreloadFunc()
	}
	return nil, nil
}

func (m *MockSuperAdminRepository) GetPaidInvoicesMonthlyRevenue(startDate, endDate *time.Time) ([]struct {
	Month        string
	Year         string
	TotalRevenue float64
}, error) {
	if m.GetPaidInvoicesMonthlyRevenueFunc != nil {
		return m.GetPaidInvoicesMonthlyRevenueFunc(startDate, endDate)
	}
	return nil, nil
}

var _ repository.SuperAdminRepository = &MockSuperAdminRepository{}

// MockRepositories holds all mock repository implementations.
type MockRepositories struct {
	AdminCompanyRepo         *MockAdminCompanyRepository
	AttendanceLocationRepo   *MockAttendanceLocationRepository
	AttendanceRepo           *MockAttendanceRepository
	BroadcastRepo            *MockBroadcastRepository
	CompanyRepo              *MockCompanyRepository
	CustomOfferRepo          *MockCustomOfferRepository
	CustomPackageRequestRepo *MockCustomPackageRequestRepository
	DivisionRepo             *MockDivisionRepository
	EmployeeRepo             *MockEmployeeRepository
	FaceImageRepo            *MockFaceImageRepository
	InvoiceRepo              *MockInvoiceRepository
	LeaveRequestRepo         *MockLeaveRequestRepository
	PasswordResetRepo        *MockPasswordResetRepository
	ShiftRepo                *MockShiftRepository
	SubscriptionPackageRepo  *MockSubscriptionPackageRepository
	SuperAdminRepo           *MockSuperAdminRepository
}

// NewMockRepositories creates and returns a new instance of MockRepositories.
func NewMockRepositories() *MockRepositories {
	return &MockRepositories{
		AdminCompanyRepo:         &MockAdminCompanyRepository{},
		AttendanceLocationRepo:   &MockAttendanceLocationRepository{},
		AttendanceRepo:           &MockAttendanceRepository{},
		BroadcastRepo:            &MockBroadcastRepository{},
		CompanyRepo:              &MockCompanyRepository{},
		CustomOfferRepo:          &MockCustomOfferRepository{},
		CustomPackageRequestRepo: &MockCustomPackageRequestRepository{},
		DivisionRepo:             &MockDivisionRepository{},
		EmployeeRepo:             &MockEmployeeRepository{},
		FaceImageRepo:            &MockFaceImageRepository{},
		InvoiceRepo:              &MockInvoiceRepository{},
				LeaveRequestRepo:         &MockLeaveRequestRepository{},
		PasswordResetRepo:        &MockPasswordResetRepository{},
		ShiftRepo:                &MockShiftRepository{},
		SubscriptionPackageRepo:  &MockSubscriptionPackageRepository{},
		SuperAdminRepo:           &MockSuperAdminRepository{},
	}
}

// Dummy variable to ensure gorm.io/gorm is imported and used
var _ *gorm.DB
