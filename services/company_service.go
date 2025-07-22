package services

import (
	"fmt"
	"go-face-auth/database"
	"go-face-auth/database/repository"
	"go-face-auth/helper"
	"go-face-auth/models"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetCompanyDetails(companyID int) (map[string]interface{}, error) {
	company, err := repository.GetCompanyByID(companyID)
	if err != nil {
		return nil, err
	}

	if company == nil {
		return nil, nil
	}

	adminCompany, err := repository.GetAdminCompanyByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	responseData := map[string]interface{}{
		"id":                    company.ID,
		"name":                  company.Name,
		"address":               company.Address,
		"admin_email":           adminCompany.Email,
		"subscription_status":   company.SubscriptionStatus,
		"trial_start_date":      company.TrialStartDate,
		"trial_end_date":        company.TrialEndDate,
		"timezone":              company.Timezone,
	}

	return responseData, nil
}

func UpdateCompanyDetails(companyID int, name, address, timezone string) (*models.CompaniesTable, error) {
	company, err := repository.GetCompanyByID(companyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, nil
	}

	if name != "" {
		company.Name = name
	}
	if address != "" {
		company.Address = address
	}
	if timezone != "" {
		_, err := time.LoadLocation(timezone)
		if err != nil {
			return nil, err
		}
		company.Timezone = timezone
	}

	if err := repository.UpdateCompany(company); err != nil {
		return nil, err
	}

	return company, nil
}

// RegisterCompanyRequest defines the structure for the company registration request body.
type RegisterCompanyRequest struct {
	CompanyName         string `json:"company_name" binding:"required"`
	CompanyAddress      string `json:"company_address"`
	AdminEmail          string `json:"admin_email" binding:"required,email"`
	AdminPassword       string `json:"admin_password" binding:"required,min=6"`
	SubscriptionPackageID int   `json:"subscription_package_id" binding:"required"`
}

func RegisterCompany(req RegisterCompanyRequest) (*models.CompaniesTable, *models.AdminCompaniesTable, error) {
	var subPackage models.SubscriptionPackageTable
	if err := database.DB.First(&subPackage, req.SubscriptionPackageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("subscription package not found")
		} else {
			return nil, nil, fmt.Errorf("failed to retrieve subscription package: %w", err)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, nil, fmt.Errorf("failed to start database transaction: %w", err)
	}

	now := time.Now()
	trialEndDate := now.AddDate(0, 0, 14)
	company := models.CompaniesTable{
		Name:                req.CompanyName,
		Address:             req.CompanyAddress,
		SubscriptionPackageID: req.SubscriptionPackageID,
		SubscriptionStatus:  "trial",
		TrialStartDate:      &now,
		TrialEndDate:        &trialEndDate,
	}
	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create company: %w", err)
	}

	confirmationToken := uuid.New().String()
	adminCompany := models.AdminCompaniesTable{
		CompanyID:         company.ID,
		Email:             req.AdminEmail,
		Password:          string(hashedPassword),
		Role:              "admin",
		ConfirmationToken: &confirmationToken,
		IsConfirmed:       false,
	}
	if err := tx.Create(&adminCompany).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create company admin: %w", err)
	}

	defaultShift := models.ShiftsTable{
		CompanyID:          company.ID,
		Name:               "Shift Pagi",
		StartTime:          "09:00:00",
		EndTime:            "17:00:00",
		GracePeriodMinutes: 15,
		IsDefault:          true,
	}
	if err := tx.Create(&defaultShift).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create default shift: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("failed to commit database transaction: %w", err)
	}

	frontendAdminBaseURL := os.Getenv("FRONTEND_ADMIN_BASE_URL")
	if frontendAdminBaseURL == "" {
		log.Println("WARNING: FRONTEND_ADMIN_BASE_URL environment variable is not set. Email confirmation link may not work correctly.")
	}
	confirmationLink := fmt.Sprintf("%s/confirm-email?token=%s", frontendAdminBaseURL, confirmationToken)

	go func() {
		err := helper.SendConfirmationEmail(req.AdminEmail, req.CompanyName, confirmationLink)
		if err != nil {
			log.Printf("Error sending confirmation email to %s: %v", req.AdminEmail, err)
		}
	}()

	return &company, &adminCompany, nil
}

func ConfirmEmail(token string) error {
	var adminCompany models.AdminCompaniesTable
	result := database.DB.Where("confirmation_token = ?", token).First(&adminCompany)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired confirmation token")
		} else {
			return fmt.Errorf("failed to retrieve admin company: %w", result.Error)
		}
	}

	if adminCompany.IsConfirmed {
		return fmt.Errorf("email already confirmed")
	}

	adminCompany.IsConfirmed = true
	adminCompany.ConfirmationToken = nil

	if err := database.DB.Save(&adminCompany).Error; err != nil {
		return fmt.Errorf("failed to confirm email: %w", err)
	}

	return nil
}