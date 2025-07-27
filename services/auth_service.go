package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
)

type AuthService interface {
	AuthenticateSuperAdmin(email, password string) (*models.SuperAdminTable, error)
	AuthenticateAdminCompany(email, password string) (*models.AdminCompaniesTable, error)
	AuthenticateEmployee(email, password string) (*models.EmployeesTable, []*models.AttendanceLocation, error)
}

type authService struct {
	superAdminRepo        repository.SuperAdminRepository
	adminCompanyRepo      repository.AdminCompanyRepository
	employeeRepo          repository.EmployeeRepository
	attendanceLocationRepo repository.AttendanceLocationRepository
}

func NewAuthService(superAdminRepo repository.SuperAdminRepository, adminCompanyRepo repository.AdminCompanyRepository, employeeRepo repository.EmployeeRepository, attendanceLocationRepo repository.AttendanceLocationRepository) AuthService {
	return &authService{
		superAdminRepo:        superAdminRepo,
		adminCompanyRepo:      adminCompanyRepo,
		employeeRepo:          employeeRepo,
		attendanceLocationRepo: attendanceLocationRepo,
	}
}

func (s *authService) AuthenticateSuperAdmin(email, password string) (*models.SuperAdminTable, error) {
	superAdmin, err := s.superAdminRepo.GetSuperAdminByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve super user: %w", err)
	}
	if superAdmin == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := helper.CheckPasswordHash(password, superAdmin.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return superAdmin, nil
}

func (s *authService) AuthenticateAdminCompany(email, password string) (*models.AdminCompaniesTable, error) {
	adminCompany, err := s.adminCompanyRepo.GetAdminCompanyByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin company: %w", err)
	}
	if adminCompany == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !adminCompany.IsConfirmed {
		return nil, fmt.Errorf("email not confirmed. Please check your inbox for a confirmation link")
	}

	if err := helper.CheckPasswordHash(password, adminCompany.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return adminCompany, nil
}

func (s *authService) AuthenticateEmployee(email, password string) (*models.EmployeesTable, []*models.AttendanceLocation, error) {
	employee, err := s.employeeRepo.GetEmployeeByEmail(email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve employee: %w", err)
	}
	if employee == nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	if err := helper.CheckPasswordHash(password, employee.Password); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	locationsValue, err := s.attendanceLocationRepo.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve company location information: %w", err)
	}

	var locations []*models.AttendanceLocation
	for i := range locationsValue {
		locations = append(locations, &locationsValue[i])
	}

	return employee, locations, nil
}