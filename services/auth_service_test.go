package services_test

import (
	"errors"
	"go-face-auth/helper"
	"go-face-auth/models"
	"go-face-auth/services"
	"testing"

	"github.com/stretchr/testify/assert"
)



	

	func TestAuthenticateSuperAdmin(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewAuthService(mocks.SuperAdminRepo, nil, nil, nil)

	password, _ := helper.HashPassword("password")
	superAdmin := &models.SuperAdminTable{Email: "super@admin.com", Password: password}

	t.Run("Success", func(t *testing.T) {
		mocks.SuperAdminRepo.GetSuperAdminByEmailFunc = func(email string) (*models.SuperAdminTable, error) {
			return superAdmin, nil
		}

		result, err := service.AuthenticateSuperAdmin("super@admin.com", "password")

		assert.NoError(t, err)
		assert.Equal(t, superAdmin, result)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		mocks.SuperAdminRepo.GetSuperAdminByEmailFunc = func(email string) (*models.SuperAdminTable, error) {
			return superAdmin, nil
		}

		_, err := service.AuthenticateSuperAdmin("super@admin.com", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Not Found", func(t *testing.T) {
		mocks.SuperAdminRepo.GetSuperAdminByEmailFunc = func(email string) (*models.SuperAdminTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.AuthenticateSuperAdmin("super@admin.com", "password")

		assert.Error(t, err)
	})
}

func TestAuthenticateAdminCompany(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewAuthService(nil, mocks.AdminCompanyRepo, nil, nil)

	password, _ := helper.HashPassword("password")
	adminCompany := &models.AdminCompaniesTable{Email: "admin@company.com", Password: password, IsConfirmed: true}

	t.Run("Success", func(t *testing.T) {
		mocks.AdminCompanyRepo.GetAdminCompanyByEmailFunc = func(email string) (*models.AdminCompaniesTable, error) {
			return adminCompany, nil
		}

		result, err := service.AuthenticateAdminCompany("admin@company.com", "password")

		assert.NoError(t, err)
		assert.Equal(t, adminCompany, result)
	})

			t.Run("Wrong Password", func(t *testing.T) {
			mocks.AdminCompanyRepo.GetAdminCompanyByEmailFunc = func(email string) (*models.AdminCompaniesTable, error) {
				return adminCompany, nil
			}

		_, err := service.AuthenticateAdminCompany("admin@company.com", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Not Confirmed", func(t *testing.T) {
		adminCompany.IsConfirmed = false
		mocks.AdminCompanyRepo.GetAdminCompanyByEmailFunc = func(email string) (*models.AdminCompaniesTable, error) {
			return adminCompany, nil
		}

		_, err := service.AuthenticateAdminCompany("admin@company.com", "password")

		assert.Error(t, err)
		assert.Equal(t, "email not confirmed. Please check your inbox for a confirmation link", err.Error())
		adminCompany.IsConfirmed = true // reset for other tests
	})

	t.Run("Not Found", func(t *testing.T) {
		mocks.AdminCompanyRepo.GetAdminCompanyByEmailFunc = func(email string) (*models.AdminCompaniesTable, error) {
			return nil, errors.New("not found")
		}

		_, err := service.AuthenticateAdminCompany("admin@company.com", "password")

		assert.Error(t, err)
	})
}

func TestAuthenticateEmployee(t *testing.T) {
	mocks := services.NewMockRepositories()
	service := services.NewAuthService(nil, nil, mocks.EmployeeRepo, mocks.AttendanceLocationRepo)

	password, _ := helper.HashPassword("password")
	employee := &models.EmployeesTable{Email: "employee@company.com", Password: password, CompanyID: 1}
	locations := []models.AttendanceLocation{{Name: "loc1"}, {Name: "loc2"}}

	t.Run("Success", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByEmailFunc = func(email string) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationsByCompanyID", uint(employee.CompanyID)).Return(locations, nil).Once()

		result, locs, err := service.AuthenticateEmployee("employee@company.com", "password")

		assert.NoError(t, err)
		assert.Equal(t, employee, result)
		assert.Equal(t, len(locations), len(locs))
	})

			t.Run("Wrong Password", func(t *testing.T) {
			mocks.EmployeeRepo.GetEmployeeByEmailFunc = func(email string) (*models.EmployeesTable, error) {
				return employee, nil
			}

		_, _, err := service.AuthenticateEmployee("employee@company.com", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Employee Not Found", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByEmailFunc = func(email string) (*models.EmployeesTable, error) {
			return nil, errors.New("not found")
		}

		_, _, err := service.AuthenticateEmployee("employee@company.com", "password")

		assert.Error(t, err)
	})

	t.Run("Location Not Found", func(t *testing.T) {
		mocks.EmployeeRepo.GetEmployeeByEmailFunc = func(email string) (*models.EmployeesTable, error) {
			return employee, nil
		}
		mocks.AttendanceLocationRepo.On("GetAttendanceLocationsByCompanyID", uint(employee.CompanyID)).Return([]models.AttendanceLocation{}, errors.New("not found")).Once()

		_, _, err := service.AuthenticateEmployee("employee@company.com", "password")

		assert.Error(t, err)
	})
}
