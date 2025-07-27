package services

import (
	"fmt"
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

type CompanyService interface {
	CreateCompany(req CreateCompanyRequest) (*models.CompaniesTable, error)
	GetCompanyByID(companyID int) (*models.CompaniesTable, error)
	GetCompanyDetails(companyID int) (map[string]interface{}, error)
	UpdateCompanyDetails(companyID int, name, address, timezone string) (*models.CompaniesTable, error)
	RegisterCompany(req RegisterCompanyRequest) (*models.CompaniesTable, *models.AdminCompaniesTable, error)
	ConfirmEmail(token string) error
	GetCompanySubscriptionStatus(companyID int) (map[string]interface{}, error)
}

type companyService struct {
	companyRepo      repository.CompanyRepository
	adminCompanyRepo repository.AdminCompanyRepository
	subscriptionRepo repository.SubscriptionPackageRepository
	shiftRepo        repository.ShiftRepository
	db               *gorm.DB
}

func NewCompanyService(companyRepo repository.CompanyRepository, adminCompanyRepo repository.AdminCompanyRepository, subscriptionRepo repository.SubscriptionPackageRepository, shiftRepo repository.ShiftRepository, db *gorm.DB) CompanyService {
	return &companyService{
		companyRepo:      companyRepo,
		adminCompanyRepo: adminCompanyRepo,
		subscriptionRepo: subscriptionRepo,
		shiftRepo:        shiftRepo,
		db:               db,
	}
}

type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func (s *companyService) CreateCompany(req CreateCompanyRequest) (*models.CompaniesTable, error) {
	company := &models.CompaniesTable{
		Name:    req.Name,
		Address: req.Address,
	}
	if err := s.companyRepo.CreateCompany(company); err != nil {
		return nil, err
	}
	return company, nil
}

func (s *companyService) GetCompanyByID(companyID int) (*models.CompaniesTable, error) {
	return s.companyRepo.GetCompanyByID(companyID)
}

func (s *companyService) GetCompanyDetails(companyID int) (map[string]interface{}, error) {
	company, err := s.companyRepo.GetCompanyByID(companyID)
	if err != nil {
		return nil, err
	}

	if company == nil {
		return nil, nil
	}

	adminCompany, err := s.adminCompanyRepo.GetAdminCompanyByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	responseData := map[string]interface{}{
		"id":                      company.ID,
		"name":                    company.Name,
		"address":                 company.Address,
		"admin_email":             adminCompany.Email,
		"subscription_status":     company.SubscriptionStatus,
		"trial_start_date":        company.TrialStartDate,
		"trial_end_date":          company.TrialEndDate,
		"timezone":                company.Timezone,
		"subscription_package_id": company.SubscriptionPackageID,
		"billing_cycle":           company.BillingCycle,
	}

	return responseData, nil
}

func (s *companyService) UpdateCompanyDetails(companyID int, name, address, timezone string) (*models.CompaniesTable, error) {
	company, err := s.companyRepo.GetCompanyByID(companyID)
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

	if err := s.companyRepo.UpdateCompany(company); err != nil {
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
	BillingCycle        string `json:"billing_cycle" binding:"required"`
}

func (s *companyService) RegisterCompany(req RegisterCompanyRequest) (*models.CompaniesTable, *models.AdminCompaniesTable, error) {
	subPackage, err := s.subscriptionRepo.GetSubscriptionPackageByID(req.SubscriptionPackageID)
	if err != nil {
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

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, nil, fmt.Errorf("failed to start database transaction: %w", err)
	}

	now := time.Now()
	trialEndDate := now.AddDate(0, 0, 14)
	company := models.CompaniesTable{
		Name:                req.CompanyName,
		Address:             req.CompanyAddress,
		SubscriptionPackageID: subPackage.ID,
		SubscriptionStatus:  "trial",
		TrialStartDate:      &now,
		TrialEndDate:        &trialEndDate,
		BillingCycle:        req.BillingCycle,
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

func (s *companyService) ConfirmEmail(token string) error {
	var adminCompany models.AdminCompaniesTable
	result := s.db.Where("confirmation_token = ?", token).First(&adminCompany)

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

	if err := s.db.Save(&adminCompany).Error; err != nil {
		return fmt.Errorf("failed to confirm email: %w", err)
	}

	return nil
}

func (s *companyService) GetCompanySubscriptionStatus(companyID int) (map[string]interface{}, error) {
	company, err := s.companyRepo.GetCompanyByID(companyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, nil
	}

	responseData := map[string]interface{}{
		"subscription_status":   company.SubscriptionStatus,
		"trial_start_date":      company.TrialStartDate,
		"trial_end_date":        company.TrialEndDate,
		"subscription_package_id": company.SubscriptionPackageID,
		"billing_cycle":         company.BillingCycle,
	}

	return responseData, nil
}