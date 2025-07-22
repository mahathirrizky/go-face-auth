package services

import (
	"fmt"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
)

func AuthenticateSuperAdmin(email, password string) (*models.SuperAdminTable, error) {
	superAdmin, err := repository.GetSuperAdminByEmail(email)
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

func AuthenticateAdminCompany(email, password string) (*models.AdminCompaniesTable, error) {
	adminCompany, err := repository.GetAdminCompanyByEmail(email)
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

func AuthenticateEmployee(email, password string) (*models.EmployeesTable, []*models.AttendanceLocation, error) {
	employee, err := repository.GetEmployeeByEmail(email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve employee: %w", err)
	}
	if employee == nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	if err := helper.CheckPasswordHash(password, employee.Password); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	locationsValue, err := repository.GetAttendanceLocationsByCompanyID(uint(employee.CompanyID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve company location information: %w", err)
	}

	var locations []*models.AttendanceLocation
	for i := range locationsValue {
		locations = append(locations, &locationsValue[i])
	}

	return employee, locations, nil
}
